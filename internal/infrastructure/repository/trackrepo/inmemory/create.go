package trackrepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/track"
)

func (r *Repository) Create(ctx context.Context, trackToCreate *track.Track) error {
	r.trackByID[trackToCreate.ID()] = trackToCreate
	url := trackToCreate.URL()
	r.trackByURL[url.String()] = trackToCreate

	return nil
}
