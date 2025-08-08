package userrepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/infrastructure/repository/userrepo"
)

func (r *Repository) GetByExternalID(ctx context.Context, externalUserID id.UserExternal) (*user.User, error) {
	targetUser, exists := r.userByExternalID[externalUserID]
	if !exists {
		return nil, userrepo.NewErrUserByExternalIDNotFound(externalUserID)
	}

	return targetUser, nil
}
