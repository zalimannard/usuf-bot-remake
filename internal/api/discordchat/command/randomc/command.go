package randomc

import (
	"context"
)

var (
	names       = []string{"random"}
	parameters  = make([]string, 0)
	description = "Случайный порядок треков"
)

type randomOrderUseCase interface {
	Random(ctx context.Context) error
}

type Command struct {
	randomOrderUseCase randomOrderUseCase
}

func New(randomOrderUseCase randomOrderUseCase) *Command {
	return &Command{
		randomOrderUseCase: randomOrderUseCase,
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
