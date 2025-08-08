package grouprepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
)

func (r *Repository) Create(ctx context.Context, groupToCreate *group.Group) error {
	r.groupByID[groupToCreate.ID()] = groupToCreate
	r.groupByExternalID[groupToCreate.ExternalID()] = groupToCreate

	return nil
}
