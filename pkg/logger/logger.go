package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config interface {
	Level() string
}

type Logger struct {
	*zerolog.Logger
}

func consoleWriter(level zerolog.Level) zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	writer.FormatTimestamp = func(i any) string {
		return ""
	}
	writer.FormatFieldName = func(i any) string {
		return fmt.Sprintf("\n    %s: ", i)
	}
	writer.FormatFieldValue = func(i any) string {
		return fmt.Sprintf("%s", i)
	}

	return writer
}

func New(cfg Config, additionalWriters ...io.Writer) (*Logger, error) {
	isLevelParsed := true
	defaultLevel := zerolog.TraceLevel

	level, err := zerolog.ParseLevel(cfg.Level())
	if err != nil {
		level = defaultLevel
		isLevelParsed = false
	}
	zerolog.SetGlobalLevel(level)

	writers := []io.Writer{
		consoleWriter(level),
	}
	writers = append(writers, additionalWriters...)

	multiWriter := io.MultiWriter(writers...)
	logger := zerolog.New(multiWriter).With().Timestamp().Logger()

	if !isLevelParsed {
		logger.Warn().Msgf("Failed to parse level: %s. Default level: %s", cfg.Level(), defaultLevel.String())
	}
	if len(additionalWriters) == 0 {
		logger.Warn().Msg("No additional writers. Only default console writer will be used")
	}

	return &Logger{
		&logger,
	}, nil
}
