package trackrepoinmemory

import (
	"context"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/infrastructure/repository/trackrepo"
)

func (r *Repository) GetByURL(ctx context.Context, targetURL url.URL) (*track.Track, error) {
	targetTrack, exists := r.trackByURL[targetURL.String()]
	if !exists {
		return nil, trackrepo.NewErrTrackByURLNotFound(targetURL)
	}

	return targetTrack, nil
}
