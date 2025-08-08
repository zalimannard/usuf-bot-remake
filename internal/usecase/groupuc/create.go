package groupuc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
)

func (u *UseCase) Create(ctx context.Context, groupToCreate *group.Group) error {
	err := u.groupProvider.Create(ctx, groupToCreate)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	return nil
}
