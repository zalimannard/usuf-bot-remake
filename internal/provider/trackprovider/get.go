package trackprovider

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/infrastructure/repository/trackrepo"
)

func (p *Provider) Get(ctx context.Context, trackID id.Track) (*track.Track, error) {
	targetTrack, err := p.trackRepository.Get(ctx, trackID)
	if err != nil {
		if errors.Is(err, trackrepo.ErrTrackNotFound) {
			return nil, fmt.Errorf("%w: %s", NewErrTrackNotFound(trackID), err.Error())
		}
		return nil, fmt.Errorf("failed to get track by url: %s", err)
	}

	return targetTrack, nil
}
