package grouprepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/infrastructure/repository/grouprepo"
)

func (r *Repository) GetByExternalID(ctx context.Context, externalGroupID id.GroupExternal) (*group.Group, error) {
	targetGroup, exists := r.groupByExternalID[externalGroupID]
	if !exists {
		return nil, grouprepo.NewErrGroupByExternalIDNotFound(externalGroupID)
	}

	return targetGroup, nil
}
