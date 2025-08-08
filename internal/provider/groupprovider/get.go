package groupprovider

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/infrastructure/repository/grouprepo"
)

func (p *Provider) Get(ctx context.Context, groupID id.Group) (*group.Group, error) {
	targetGroup, err := p.groupRepository.Get(ctx, groupID)
	if err != nil {
		if errors.Is(err, grouprepo.ErrGroupNotFound) {
			return nil, newErrGroupNotFound(groupID)
		}
		return nil, fmt.Errorf("failed to get group in repository: %s", err.Error())
	}

	return targetGroup, nil
}
