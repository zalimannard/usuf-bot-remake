package queueprovider

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/infrastructure/repository/queuerepo"
)

func (p *Provider) GetByGroupID(ctx context.Context, groupID id.Group) (*queue.Queue, error) {
	targetQueue, err := p.queueRepository.GetByGroupID(ctx, groupID)
	if err == nil {
		return targetQueue, nil
	}

	if !errors.Is(err, queuerepo.ErrQueueByGroupIDNotFound) {
		return nil, fmt.Errorf("failed to get queue by group id: %s", err)
	}

	newQueue, err := queue.New(
		nil,
		make([]queue.Item, 0),
		queue.OrderTypeNormal,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct new queue: %s", err)
	}

	err = p.queueRepository.Create(ctx, groupID, newQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to create new queue: %s", err)
	}

	return newQueue, nil
}
