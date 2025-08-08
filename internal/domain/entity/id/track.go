package id

import (
	"fmt"

	"github.com/google/uuid"
)

type Track uuid.UUID

func ParseTrack(id string) (Track, error) {
	trackID, err := uuid.Parse(id)
	if err != nil {
		return Track(trackID), fmt.Errorf("invalid track id: %w", err)
	}
	return Track(trackID), nil
}

func GenerateTrack() Track {
	return Track(uuid.New())
}

func (g Track) String() string {
	return uuid.UUID(g).String()
}
