package trackrepoinmemory

import (
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/track"
)

type Repository struct {
	trackByID  map[id.Track]*track.Track
	trackByURL map[string]*track.Track
}

func New() *Repository {
	return &Repository{
		trackByID:  make(map[id.Track]*track.Track),
		trackByURL: make(map[string]*track.Track),
	}
}
