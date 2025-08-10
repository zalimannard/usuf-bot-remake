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
	attemptTimeout = 30 * time.Second // было 10s — на холодном старте может не хватать
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
	var cancelPrev context.CancelFunc

	for attempt := 0; attempt < attemptsCount; attempt++ {
		fmt.Printf("Attempt %d...\n", attempt+1)

		// Отменяем предыдущую попытку (если была)
		if cancelPrev != nil {
			cancelPrev()
		}
		attemptCtx, cancel := context.WithCancel(ctx)
		cancelPrev = cancel

		opusCh, errCh := play(attemptCtx, targetURL)

		timer := time.NewTimer(attemptTimeout)
		select {
		case <-ctx.Done():
			timer.Stop()
			cancel()
			return nil, nil

		case <-timer.C:
			// Не дождались первого пакета — мягко отменяем попытку и пробуем снова
			timer.Stop()
			cancel()
			drainErrCh(errCh, 2*time.Second)
			continue

		case err, ok := <-errCh:
			// Получили ошибку до первого пакета
			if !ok {
				// попытка завершилась без деталей — пробуем ещё
				cancel()
				continue
			}
			if errors.Is(err, ErrEndOfStream) {
				// Конец потока до первого пакета — возвратим пустой opus-канал, errCh пробросим вызывающему
				empty := make(chan []byte)
				close(empty)
				return empty, errCh
			}
			fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
			// Отменяем и пробуем следующую
			cancel()
			drainErrCh(errCh, 2*time.Second)
			continue

		case first, ok := <-opusCh:
			// Первый opus-пакет получен — успех
			if !ok {
				// канал закрылся слишком рано — пробуем снова
				cancel()
				continue
			}
			timer.Stop()

			// Раз мы успешно стартовали — больше не отменяем эту попытку
			cancelPrev = nil

			out := make(chan []byte, 16)
			out <- first

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

	// Все попытки исчерпаны
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
			"-vn", // видео не нужно
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

		// Закрыть pw после завершения yt-dlp (чтобы ffmpeg получил EOF)
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

		// Ждём ffmpeg в отдельной горутине
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
				// нормальное завершение потока
				if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
					strings.Contains(err.Error(), "file already closed") {
					if !wroteAny {
						// Сообщим наверх, что поток завершился до первого пакета
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
				wroteAny = true
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

// best-effort дренаж канала ошибок, чтобы дать горутине корректно завершиться
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
