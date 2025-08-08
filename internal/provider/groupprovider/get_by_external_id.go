package groupprovider

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/infrastructure/repository/grouprepo"
)

func (p *Provider) GetByExternalID(ctx context.Context, externalGroupID id.GroupExternal) (*group.Group, error) {
	targetGroup, err := p.groupRepository.GetByExternalID(ctx, externalGroupID)
	if err != nil {
		if errors.Is(err, grouprepo.ErrGroupByExternalIDNotFound) {
			return nil, newErrGroupByExternalIDNotFound(externalGroupID)
		}
		return nil, fmt.Errorf("failed to get group by external id in repository: %s", err.Error())
	}

	return targetGroup, nil
}
