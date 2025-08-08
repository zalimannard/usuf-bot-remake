package playc

import (
	"context"
	"net/url"
)

var (
	names       = []string{"play", "p"}
	parameters  = []string{"URL", "Запрос"}
	description = "Добавить трек(и) в конец очереди"
)

type playUseCase interface {
	PlayByURLs(ctx context.Context, urls []url.URL) error
	PlayByQuery(ctx context.Context, query string) error
}

type Command struct {
	playUseCase playUseCase
}

func New(playUseCase playUseCase) *Command {
	return &Command{
		playUseCase: playUseCase,
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
