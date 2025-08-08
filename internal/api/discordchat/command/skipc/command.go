package skipc

import (
	"context"
	"errors"
)

var (
	names       = []string{"skip", "s"}
	parameters  = make([]string, 0)
	description = "Пропустить текущий трек"

	ErrFailedToSkipTrack = errors.New("failed to skip track")
)

type skipUseCase interface {
	Skip(ctx context.Context) error
}

type Command struct {
	skipUseCase skipUseCase
}

func New(skipUseCase skipUseCase) *Command {
	return &Command{
		skipUseCase: skipUseCase,
	}
}

func (c *Command) Names() []string {
	return names
}

func (c *Command) Parameters() []string {
	return parameters
}

func (c *Command) Description() string {
	return description
}
