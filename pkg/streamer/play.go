package streamer

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"layeh.com/gopus"
)

// Мощный вайбкодинг

const (
	attemptsCount  = 5
	attemptTimeout = 30 * time.Second
	frameSize      = 960
	sampleRate     = 48000
	channels       = 2

	// Если потребитель перестал читать — не висим бесконечно.
	sendBlockTimeout = 3 * time.Second
)

var (
	ErrEndOfStream         = errors.New("end of stream")
	ErrPlaybackUnavailable = errors.New("playback unavailable")
	ErrConsumerBlocked     = errors.New("consumer is not reading (blocked)")
)

// Play пробует запустить play до attemptsCount раз.
// Возвращает готовые каналы: opus (payloads) и errors.
func Play(ctx context.Context, targetURL url.URL) (<-chan []byte, <-chan error) {
	var cancelPrev context.CancelFunc

	for attempt := 0; attempt < attemptsCount; attempt++ {
		fmt.Printf("Attempt %d...\n", attempt+1)

		if cancelPrev != nil {
			cancelPrev()
		}
		attemptCtx, cancel := context.WithCancel(ctx)
		cancelPrev = cancel

		opusCh, errCh := play(attemptCtx, targetURL)

		timer := time.NewTimer(attemptTimeout)
		select {
		case <-ctx.Done():
			fmt.Println("1")
			timer.Stop()
			cancel()
			return nil, nil

		case <-timer.C:
			fmt.Println("2")
			timer.Stop()
			cancel()
			drainErrCh(errCh, 2*time.Second)
			continue

		case err, ok := <-errCh:
			fmt.Println("3")
			if !ok {
				cancel()
				continue
			}
			if errors.Is(err, ErrEndOfStream) {
				empty := make(chan []byte)
				close(empty)
				return empty, errCh
			}
			fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
			cancel()
			drainErrCh(errCh, 2*time.Second)
			continue

		case first, ok := <-opusCh:
			fmt.Println("4")
			if !ok {
				cancel()
				continue
			}
			timer.Stop()
			cancelPrev = nil

			out := make(chan []byte, 16)
			out <- first

			// Форвардер opusCh -> out с тайм-аутом на случай, если потребитель завис.
			go func() {
				defer close(out)
				for {
					select {
					case <-ctx.Done():
						return
					case p, ok := <-opusCh:
						if !ok {
							return
						}
						select {
						case out <- p:
						case <-ctx.Done():
							return
						case <-time.After(sendBlockTimeout):
							// Потребитель не читает — прекращаем, чтобы не зависнуть
							return
						}
					}
				}
			}()
			return out, errCh
		}
	}

	if cancelPrev != nil {
		cancelPrev()
	}
	opusCh := make(chan []byte, 1)
	errCh := make(chan error, 1)
	errCh <- ErrPlaybackUnavailable
	return opusCh, errCh
}

func play(ctx context.Context, targetURL url.URL) (chan []byte, chan error) {
	opusChan := make(chan []byte, 16)
	errChan := make(chan error, 2)

	go func() {
		var wg sync.WaitGroup
		defer func() {
			wg.Wait()
			close(errChan)
		}()
		defer close(opusChan)

		pr, pw := io.Pipe()

		ytCmd := exec.CommandContext(ctx, "yt-dlp", "-o", "-", targetURL.String())
		ytCmd.Stdout = pw
		ytCmd.Stderr = os.Stderr

		ffmpegCmd := exec.CommandContext(ctx,
			"ffmpeg",
			"-nostdin", "-hide_banner",
			"-i", "pipe:0",
			"-vn",
			"-f", "s16le",
			"-ar", fmt.Sprint(sampleRate),
			"-ac", fmt.Sprint(channels),
			"pipe:1",
		)
		ffmpegCmd.Stdin = pr
		ffmpegCmd.Stderr = os.Stderr

		ffOut, err := ffmpegCmd.StdoutPipe()
		if err != nil {
			select {
			case errChan <- fmt.Errorf("ffmpeg stdout pipe: %w", err):
			default:
			}
			_ = pw.Close()
			_ = pr.Close()
			return
		}

		if err := ffmpegCmd.Start(); err != nil {
			select {
			case errChan <- fmt.Errorf("failed to start ffmpeg: %w", err):
			default:
			}
			_ = pw.Close()
			_ = pr.Close()
			return
		}

		if err := ytCmd.Start(); err != nil {
			select {
			case errChan <- fmt.Errorf("failed to start yt-dlp: %w", err):
			default:
			}
			_ = tryKill(ffmpegCmd)
			_ = ffmpegCmd.Wait()
			_ = pw.Close()
			_ = pr.Close()
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := ytCmd.Wait()
			_ = pw.Close()
			if err != nil {
				select {
				case errChan <- fmt.Errorf("yt-dlp wait error: %w", err):
				default:
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ffmpegCmd.Wait()
		}()

		encoder, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
		if err != nil {
			select {
			case errChan <- fmt.Errorf("failed to create encoder: %w", err):
			default:
			}
			_ = tryKill(ytCmd, ffmpegCmd)
			_ = pw.Close()
			_ = pr.Close()
			return
		}

		reader := bufio.NewReader(ffOut)
		wroteAny := false

		for {
			select {
			case <-ctx.Done():
				_ = tryKill(ytCmd, ffmpegCmd)
				_ = pw.Close()
				_ = pr.Close()
				return
			default:
			}

			pcm := make([]int16, frameSize*channels)
			if err := binary.Read(reader, binary.LittleEndian, pcm); err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
					strings.Contains(err.Error(), "file already closed") {
					if !wroteAny {
						select {
						case errChan <- ErrEndOfStream:
						default:
						}
					}
					_ = pw.Close()
					_ = pr.Close()
					_ = tryKill(ytCmd, ffmpegCmd)
					return
				}

				select {
				case errChan <- fmt.Errorf("failed to read from ffmpeg: %w", err):
				default:
				}
				_ = pw.Close()
				_ = pr.Close()
				_ = tryKill(ytCmd, ffmpegCmd)
				return
			}

			opus, err := encoder.Encode(pcm, frameSize, len(pcm)*2)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("failed to encode: %w", err):
				default:
				}
				_ = pw.Close()
				_ = pr.Close()
				_ = tryKill(ytCmd, ffmpegCmd)
				return
			}

			// Критично: не блокируемся бесконечно, если нас перестали читать.
			select {
			case <-ctx.Done():
				_ = tryKill(ytCmd, ffmpegCmd)
				_ = pw.Close()
				_ = pr.Close()
				return
			case opusChan <- opus:
				wroteAny = true
			case <-time.After(sendBlockTimeout):
				select {
				case errChan <- ErrConsumerBlocked:
				default:
				}
				_ = pw.Close()
				_ = pr.Close()
				_ = tryKill(ytCmd, ffmpegCmd)
				return
			}
		}
	}()

	return opusChan, errChan
}

func tryKill(cmds ...*exec.Cmd) error {
	for _, c := range cmds {
		if c == nil || c.Process == nil {
			continue
		}
		_ = c.Process.Kill()
	}
	return nil
}

func drainErrCh(ch <-chan error, d time.Duration) {
	if ch == nil {
		return
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-timer.C:
			return
		}
	}
}
