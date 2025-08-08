package queuerepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
)

func (r *Repository) Create(ctx context.Context, groupID id.Group, newQueue *queue.Queue) error {
	r.queueByGroupID[groupID] = newQueue
	r.queueByID[newQueue.ID()] = newQueue
	r.groupIDByQueueID[newQueue.ID()] = groupID
	r.queueIDByGroupID[groupID] = newQueue.ID()

	return nil
}
