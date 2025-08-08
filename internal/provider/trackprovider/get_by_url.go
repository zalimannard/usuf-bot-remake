package trackprovider

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/infrastructure/repository/trackrepo"
)

func (p *Provider) GetByURL(ctx context.Context, targetURL url.URL) (*track.Track, error) {
	targetTrack, err := p.trackRepository.GetByURL(ctx, targetURL)
	if err == nil {
		return targetTrack, nil
	}

	if !errors.Is(err, trackrepo.ErrTrackNotFound) {
		return nil, fmt.Errorf("failed to get track by url: %s", err)
	}

	loadedTrack, err := p.trackLoader.Load(ctx, targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to load track: %s", err)
	}

	err = p.trackRepository.Create(ctx, loadedTrack)
	if err != nil {
		return nil, fmt.Errorf("failed to create track: %s", err)
	}

	return loadedTrack, nil
}
