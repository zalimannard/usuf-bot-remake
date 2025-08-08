package queuerepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/infrastructure/repository/queuerepo"
)

func (r *Repository) GetByGroupID(ctx context.Context, groupID id.Group) (*queue.Queue, error) {
	targetQueue, exists := r.queueByGroupID[groupID]
	if !exists {
		return nil, queuerepo.NewErrQueueByGroupIDNotFound(groupID)
	}

	return targetQueue, nil
}
