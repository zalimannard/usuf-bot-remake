package trackrepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/infrastructure/repository/trackrepo"
)

func (r *Repository) Get(ctx context.Context, trackID id.Track) (*track.Track, error) {
	targetTrack, exists := r.trackByID[trackID]
	if !exists {
		return nil, trackrepo.NewErrTrackNotFound(trackID)
	}

	return targetTrack, nil
}
