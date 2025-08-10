package trackprovider

import (
	"context"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/track"
)

type Provider struct {
	trackRepository trackRepository
	trackLoader     trackLoader
	infoRequester   infoRequester
}

func New(trackRepository trackRepository, trackLoader trackLoader, infoRequester infoRequester) *Provider {
	return &Provider{
		trackRepository: trackRepository,
		trackLoader:     trackLoader,
		infoRequester:   infoRequester,
	}
}

type trackRepository interface {
	Create(ctx context.Context, trackToCreate *track.Track) error
	GetByURL(ctx context.Context, targetURL url.URL) (*track.Track, error)
	Get(ctx context.Context, trackID id.Track) (*track.Track, error)
}

type trackLoader interface {
	Load(ctx context.Context, targetURL url.URL) (*track.Track, error)
}

type infoRequester interface {
	PlaylistURLs(ctx context.Context, playlistID string) ([]url.URL, error)
	GetURLByQuery(ctx context.Context, query string) (*url.URL, error)
}
