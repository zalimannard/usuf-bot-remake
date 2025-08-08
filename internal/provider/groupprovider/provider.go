package groupprovider

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
)

type Provider struct {
	groupRepository groupRepository
}

func New(groupRepository groupRepository) *Provider {
	return &Provider{
		groupRepository: groupRepository,
	}
}

type groupRepository interface {
	Create(ctx context.Context, groupToCreate *group.Group) error
	Get(ctx context.Context, groupID id.Group) (*group.Group, error)
	GetByExternalID(ctx context.Context, externalGroupID id.GroupExternal) (*group.Group, error)
}
