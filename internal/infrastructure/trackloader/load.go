package trackloader

import (
	"context"
	"net/url"
	"time"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/util"
)

func (l *Loader) Load(ctx context.Context, targetURL url.URL) (*track.Track, error) {
	return track.New(
		nil,
		"Без названия",
		util.Ptr("Неизвестен"),
		util.Ptr(5*time.Minute),
		targetURL,
		nil,
	), nil
}
