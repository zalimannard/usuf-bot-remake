package grouprepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/infrastructure/repository/grouprepo"
)

func (r *Repository) Get(ctx context.Context, groupID id.Group) (*group.Group, error) {
	targetGroup, exists := r.groupByID[groupID]
	if !exists {
		return nil, grouprepo.NewErrGroupNotFound(groupID)
	}

	return targetGroup, nil
}
