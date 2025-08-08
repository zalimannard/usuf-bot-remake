package groupuc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
)

func (u *UseCase) GetByExternalID(ctx context.Context, externalGroupID id.GroupExternal) (*group.Group, error) {
	targetGroup, err := u.groupProvider.GetByExternalID(ctx, externalGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group by external id: %w", err)
	}

	return targetGroup, nil
}
