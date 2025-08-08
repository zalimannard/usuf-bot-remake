package loopqc

import (
	"context"
)

var (
	names       = []string{"loopq", "lq"}
	parameters  = make([]string, 0)
	description = "Зациклить очередь"
)

type loopqUseCase interface {
	Loopq(ctx context.Context) error
}

type Command struct {
	loopqUseCase loopqUseCase
}

func New(loopqUseCase loopqUseCase) *Command {
	return &Command{
		loopqUseCase: loopqUseCase,
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
