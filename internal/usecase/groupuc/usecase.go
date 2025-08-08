package groupuc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
)

type groupProvider interface {
	Create(ctx context.Context, groupToCreate *group.Group) error
	GetByExternalID(ctx context.Context, externalGroupID id.GroupExternal) (*group.Group, error)
}

type UseCase struct {
	groupProvider groupProvider
}

func New(groupProvider groupProvider) *UseCase {
	return &UseCase{
		groupProvider: groupProvider,
	}
}
