package queueprovider

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

func (p *Provider) Delete(ctx context.Context, queueID id.Queue) error {
	err := p.queueRepository.Delete(ctx, queueID)
	if err != nil {
		return fmt.Errorf("failed to delete queue in repository: %s", err)
	}

	return nil
}
