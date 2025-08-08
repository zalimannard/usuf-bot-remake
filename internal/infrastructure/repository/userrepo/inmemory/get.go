package userrepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/infrastructure/repository/userrepo"
)

func (r *Repository) Get(ctx context.Context, userID id.User) (*user.User, error) {
	targetUser, exists := r.userByID[userID]
	if !exists {
		return nil, userrepo.NewErrUserNotFound(userID)
	}

	return targetUser, nil
}
