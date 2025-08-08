package playc

import (
	"context"
	"errors"
	"net/url"
)

var (
	names       = []string{"play", "p"}
	parameters  = []string{"URL", "Запрос"}
	description = "Добавить трек(и) в конец очереди"

	ErrFailedToPlayByURLs  = errors.New("failed to play by urls")
	ErrFailedToPlayByQuery = errors.New("failed to play by query")
)

type playUseCase interface {
	PlayByURLs(ctx context.Context, urls []url.URL) error
	PlayByQuery(ctx context.Context, query string) error
}

type alertUseCase interface {
	SendError(err error)
}

type Command struct {
	playUseCase  playUseCase
	alertUseCase alertUseCase
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
