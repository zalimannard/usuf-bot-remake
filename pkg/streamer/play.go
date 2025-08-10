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
	"time"

	"layeh.com/gopus"
)

// Мощный вайбкодинг

const (
	attemptsCount  = 5
	attemptTimeout = 10 * time.Second
	frameSize      = 960
	sampleRate     = 48000
	channels       = 2
)

var (
	ErrEndOfStream         = errors.New("end of stream")
	ErrPlaybackUnavailable = errors.New("playback unavailable")
)

// Play пробует запустить play до attemptsCount раз.
// Возвращает готовые каналы: opus (payloads) и errors.
func Play(ctx context.Context, targetURL url.URL) (<-chan []byte, <-chan error) {
	for attempt := 0; attempt < attemptsCount; attempt++ {
		fmt.Printf("Attempt %d...\n", attempt+1)
		opusCh, errCh := play(ctx, targetURL)

		timer := time.NewTimer(attemptTimeout)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, nil
		case <-timer.C:
			timer.Stop()
			// таймаут — попробуем снова
			continue
		case err, ok := <-errCh:
			// получили ошибку до первого пакета
			if !ok {
				// канал ошибок закрыт — пробуем заново
				continue
			}
			// Если поток просто закончился — не пытаемся перезапускать бесконечно.
			if errors.Is(err, ErrEndOfStream) {
				// Вернём канал ошибок и пустой opus-канал (закрытый) — caller увидит конец.
				empty := make(chan []byte)
				close(empty)
				return empty, errCh
			}
			fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
			// иначе пробуем следующую попытку
			continue
		case first, ok := <-opusCh:
			// получили первый opus-пакет — нужно вернуть канал, но не потерять пакет
			if !ok {
				// канал закрыт — попробуем заново
				continue
			}
			out := make(chan []byte, 16)
			// положим первый пакет
			out <- first
			// форвардер — перекладывает оставшиеся пакеты из opusCh в out
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
						}
					}
				}
			}()
			return out, errCh
		}
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
		defer close(errChan)
		defer close(opusChan)

		pr, pw := io.Pipe()

		ytCmd := exec.CommandContext(ctx, "yt-dlp", "-o", "-", targetURL.String())
		ytCmd.Stdout = pw
		ytCmd.Stderr = os.Stderr

		ffmpegCmd := exec.CommandContext(ctx,
			"ffmpeg",
			"-nostdin", "-hide_banner", // уменьшает шум в логах
			"-i", "pipe:0",
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
			_ = ffmpegCmd.Process.Kill()
			ffmpegCmd.Wait()
			_ = pw.Close()
			_ = pr.Close()
			return
		}

		// Закрыть pw после завершения yt-dlp (чтобы ffmpeg получил EOF)
		go func() {
			err := ytCmd.Wait()
			_ = pw.Close()
			if err != nil {
				select {
				case errChan <- fmt.Errorf("yt-dlp wait error: %w", err):
				default:
				}
			}
		}()

		// Ждём ffmpeg в отдельной горутине — но не блокируем main loop
		go func() {
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

		for {
			// cancellation check
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
				// нормальное завершение потока
				if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
					strings.Contains(err.Error(), "file already closed") {
					// нормальное завершение — не считать это ошибкой
					_ = pw.Close()
					_ = pr.Close()
					_ = tryKill(ytCmd, ffmpegCmd)
					return // opusChan закроется defer'ом
				}

				// прочая ошибка чтения
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

			select {
			case <-ctx.Done():
				_ = tryKill(ytCmd, ffmpegCmd)
				_ = pw.Close()
				_ = pr.Close()
				return
			case opusChan <- opus:
			}
		}
	}()

	return opusChan, errChan
}

// вспомогательная функция — пытается аккуратно убить процессы, игнорируя ошибки
func tryKill(cmds ...*exec.Cmd) error {
	for _, c := range cmds {
		if c == nil || c.Process == nil {
			continue
		}
		_ = c.Process.Kill()
	}
	return nil
}
