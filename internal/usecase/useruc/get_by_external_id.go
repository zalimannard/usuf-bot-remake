package useruc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
)

func (u *UseCase) GetByExternalID(ctx context.Context, externalUserID id.UserExternal) (*user.User, error) {
	targetUser, err := u.userProvider.GetByExternalID(ctx, externalUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by external id: %w", err)
	}

	return targetUser, nil
}
