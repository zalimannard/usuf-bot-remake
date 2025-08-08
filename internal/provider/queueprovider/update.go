package queueprovider

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/queue"
)

func (p *Provider) Update(ctx context.Context, queueToUpdate *queue.Queue) error {
	err := p.queueRepository.Update(ctx, queueToUpdate)
	if err != nil {
		return fmt.Errorf("failed to update queue in repository: %s", err)
	}

	return nil
}
