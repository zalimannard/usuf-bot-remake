package clearc

import (
	"context"
)

var (
	names       = []string{"clear", "c"}
	parameters  = make([]string, 0)
	description = "Очистить очередь"
)

type clearUseCase interface {
	Clear(ctx context.Context) error
}

type Command struct {
	clearUseCase clearUseCase
}

func New(clearUseCase clearUseCase) *Command {
	return &Command{
		clearUseCase: clearUseCase,
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
