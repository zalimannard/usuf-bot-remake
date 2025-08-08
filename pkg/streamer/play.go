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

func Play(ctx context.Context, targetURL url.URL) (<-chan []byte, <-chan error) {
	var opusChan chan []byte
	var errChan chan error

	for attempt := range attemptsCount {
		fmt.Printf("Attempt %d...\n", attempt+1)
		opusChan, errChan = play(ctx, targetURL)

		timer := time.NewTimer(attemptTimeout)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, nil
		case <-timer.C:
			timer.Stop()
			continue
		case opus := <-opusChan:
			opusChan <- opus
			return opusChan, errChan
		case err := <-errChan:
			fmt.Printf("Attempt %d failed: %s\n", attempt+1, err.Error())
		}
	}

	opusChan = make(chan []byte, 1)
	errChan = make(chan error, 1)
	errChan <- ErrPlaybackUnavailable

	return opusChan, errChan
}

func play(ctx context.Context, targetURL url.URL) (chan []byte, chan error) {
	opusChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func(opusChan chan []byte, errChan chan error) {
		defer close(errChan)
		defer close(opusChan)

		commandYtDlp := exec.Command("yt-dlp",
			"-o", "-", targetURL.String())
		commandFfmpeg := exec.Command("ffmpeg",
			"-i", "pipe:0",
			"-f", "s16le",
			"-ar", fmt.Sprint(sampleRate),
			"-ac", fmt.Sprint(channels),
			"pipe:1",
		)

		outYtDlp, err := commandYtDlp.StdoutPipe()
		if err != nil {
			errChan <- fmt.Errorf("failed to get stdout pipe for yt-dlp: %w", err)
			return
		}
		defer outYtDlp.Close()

		commandFfmpeg.Stdin = outYtDlp

		outFfmpeg, err := commandFfmpeg.StdoutPipe()
		if err != nil {
			errChan <- fmt.Errorf("failed to get stdout pipe for ffmpeg: %w", err)
			return
		}
		defer outFfmpeg.Close()

		commandYtDlp.Stderr = os.Stderr
		commandFfmpeg.Stderr = os.Stderr

		err = commandYtDlp.Start()
		if err != nil {
			errChan <- fmt.Errorf("failed to start yt-dlp: %w", err)
		}

		err = commandFfmpeg.Start()
		if err != nil {
			errChan <- fmt.Errorf("failed to start ffmpeg: %w", err)
		}

		reader := bufio.NewReader(outFfmpeg)
		encoder, err := gopus.NewEncoder(sampleRate, channels, gopus.Audio)
		if err != nil {
			errChan <- fmt.Errorf("failed to create encoder: %w", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			pcm := make([]int16, frameSize*channels)
			err = binary.Read(reader, binary.LittleEndian, &pcm)
			if err != nil {
				if errors.Is(err, io.ErrUnexpectedEOF) {
					errChan <- ErrEndOfStream
				}
				errChan <- fmt.Errorf("failed to read from ffmpeg: %w", err)
				return
			}

			opus, err := encoder.Encode(pcm, frameSize, len(pcm)*2)
			if err != nil {
				errChan <- fmt.Errorf("failed to encode: %w", err)
				return
			}

			select {
			case <-ctx.Done():
				return
			case opusChan <- opus:
			}
		}
	}(opusChan, errChan)

	return opusChan, errChan
}
