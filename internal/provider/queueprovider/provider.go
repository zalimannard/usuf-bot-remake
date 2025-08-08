package queueprovider

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
)

type Provider struct {
	queueRepository queueRepository
}

func New(queueRepository queueRepository) *Provider {
	return &Provider{
		queueRepository: queueRepository,
	}
}

type queueRepository interface {
	Create(ctx context.Context, groupID id.Group, newQueue *queue.Queue) error
	Update(ctx context.Context, queueToUpdate *queue.Queue) error
	Delete(ctx context.Context, queueID id.Queue) error
	GetByGroupID(ctx context.Context, groupID id.Group) (*queue.Queue, error)
}
