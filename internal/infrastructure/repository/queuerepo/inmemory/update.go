package queuerepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/infrastructure/repository/queuerepo"
)

func (r *Repository) Update(ctx context.Context, queueToUpdate *queue.Queue) error {
	groupID, exists := r.groupIDByQueueID[queueToUpdate.ID()]
	if !exists {
		return queuerepo.NewErrQueueNotFound(queueToUpdate.ID())
	}

	r.queueByGroupID[groupID] = queueToUpdate
	r.queueByID[queueToUpdate.ID()] = queueToUpdate
	r.groupIDByQueueID[queueToUpdate.ID()] = groupID
	r.queueIDByGroupID[groupID] = queueToUpdate.ID()

	return nil
}
