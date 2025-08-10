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
			// получили ошибку сразу — логируем и пробуем ещё
			if ok {
				fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
			}
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

		// используем io.Pipe чтобы надёжно связать stdout yt-dlp с stdin ffmpeg
		pr, pw := io.Pipe()

		ytCmd := exec.CommandContext(ctx, "yt-dlp", "-o", "-", targetURL.String())
		ytCmd.Stdout = pw
		ytCmd.Stderr = os.Stderr

		ffmpegCmd := exec.CommandContext(ctx,
			"ffmpeg",
			"-i", "pipe:0",
			"-f", "s16le",
			"-ar", fmt.Sprint(sampleRate),
			"-ac", fmt.Sprint(channels),
			"pipe:1",
		)
		ffmpegCmd.Stdin = pr
		ffmpegCmd.Stderr = os.Stderr

		// получаем stdout ffmpeg (PCM)
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

		// start ffmpeg first, чтобы он сразу был готов читать stdin
		if err := ffmpegCmd.Start(); err != nil {
			select {
			case errChan <- fmt.Errorf("failed to start ffmpeg: %w", err):
			default:
			}
			_ = pw.Close()
			_ = pr.Close()
			return
		}

		// start yt-dlp (пишет в pw)
		if err := ytCmd.Start(); err != nil {
			select {
			case errChan <- fmt.Errorf("failed to start yt-dlp: %w", err):
			default:
			}
			// kill ffmpeg if yt-dlp не запустился
			_ = ffmpegCmd.Process.Kill()
			ffmpegCmd.Wait()
			_ = pw.Close()
			_ = pr.Close()
			return
		}

		// после завершения yt-dlp нужно закрыть writer, чтобы ffmpeg получил EOF
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

		// отдельно ждём ffmpeg, но не блокируем основной поток — дождёмся после основного чтения
		ffmpegDone := make(chan error, 1)
		go func() {
			err := ffmpegCmd.Wait()
			ffmpegDone <- err
			close(ffmpegDone)
		}()

		// создаём энкодер opus
		encoder, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
		if err != nil {
			select {
			case errChan <- fmt.Errorf("failed to create encoder: %w", err):
			default:
			}
			// пытаемся корректно завершить процессы
			_ = ytCmd.Process.Kill()
			_ = ffmpegCmd.Process.Kill()
			<-ffmpegDone
			return
		}

		reader := bufio.NewReader(ffOut)

		for {
			// проверяем cancellation
			select {
			case <-ctx.Done():
				// контекст отменён — корректно убиваем процессы (CommandContext уже должен сделать это,
				// но на всякий случай)
				if ytCmd.Process != nil {
					_ = ytCmd.Process.Kill()
				}
				if ffmpegCmd.Process != nil {
					_ = ffmpegCmd.Process.Kill()
				}
				// дождёмся завершения ffmpeg
				<-ffmpegDone
				return
			default:
			}

			pcm := make([]int16, frameSize*channels)
			if err := binary.Read(reader, binary.LittleEndian, pcm); err != nil {
				fmt.Println("LOOK AT ME", err, err.Error())
				if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), "file already closed") {
					// конец потока — корректно завершаем
					select {
					case errChan <- ErrEndOfStream:
					default:
					}
					<-ffmpegDone
					return
				}
				select {
				case errChan <- fmt.Errorf("failed to read from ffmpeg: %w", err):
				default:
				}
				<-ffmpegDone
				return
			}

			opus, err := encoder.Encode(pcm, frameSize, len(pcm)*2)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("failed to encode: %w", err):
				default:
				}
				<-ffmpegDone
				return
			}

			// безопасная отправка (если ctx done — выходим)
			select {
			case <-ctx.Done():
				if ytCmd.Process != nil {
					_ = ytCmd.Process.Kill()
				}
				if ffmpegCmd.Process != nil {
					_ = ffmpegCmd.Process.Kill()
				}
				<-ffmpegDone
				return
			case opusChan <- opus:
			}
		}
	}()

	return opusChan, errChan
}
