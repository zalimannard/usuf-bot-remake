package queuerepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/infrastructure/repository/queuerepo"
)

func (r *Repository) Delete(ctx context.Context, queueID id.Queue) error {
	groupID, exists := r.groupIDByQueueID[queueID]
	if !exists {
		return queuerepo.NewErrQueueNotFound(queueID)
	}

	delete(r.queueByID, queueID)
	delete(r.groupIDByQueueID, queueID)
	delete(r.queueByGroupID, groupID)
	delete(r.queueIDByGroupID, groupID)

	return nil
}
