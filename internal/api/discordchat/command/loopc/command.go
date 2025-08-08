package loopc

import (
	"context"
)

var (
	names       = []string{"loop", "l"}
	parameters  = make([]string, 0)
	description = "Зациклить текущий трек"
)

type loopUseCase interface {
	Loop(ctx context.Context) error
}

type Command struct {
	loopUseCase loopUseCase
}

func New(loopUseCase loopUseCase) *Command {
	return &Command{
		loopUseCase: loopUseCase,
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
