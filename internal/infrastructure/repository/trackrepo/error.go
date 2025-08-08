package trackrepo

import (
	"errors"
	"fmt"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrTrackNotFound      = errors.New("track  not found")
	ErrTrackByURLNotFound = errors.New("track  not found")
)

func NewErrTrackNotFound(trackID id.Track) error {
	return fmt.Errorf("%w (id=%s)", ErrTrackNotFound, trackID)
}

func NewErrTrackByURLNotFound(targetURL url.URL) error {
	return fmt.Errorf("%w (url=%s)", ErrTrackNotFound, targetURL.String())
}
