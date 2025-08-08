package helpc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/helprow"
)

var (
	names       = []string{"help", "h"}
	parameters  = make([]string, 0)
	description = "Открыть это меню"
)

type helpUseCase interface {
	Help(ctx context.Context, rows []helprow.Row) error
}

type Command struct {
	helpUseCase helpUseCase
}

func New(helpUseCase helpUseCase) *Command {
	return &Command{
		helpUseCase: helpUseCase,
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
