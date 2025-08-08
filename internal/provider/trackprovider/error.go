package trackprovider

import (
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrTrackNotFound = errors.New("track  not found")
)

func NewErrTrackNotFound(trackID id.Track) error {
	return fmt.Errorf("%w (id=%s)", ErrTrackNotFound, trackID)
}
