package useruc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
)

type userProvider interface {
	Create(ctx context.Context, userToCreate *user.User) error
	Update(ctx context.Context, userToUpdate *user.User) error
	GetByExternalID(ctx context.Context, externalUserID id.UserExternal) (*user.User, error)
}

type UseCase struct {
	userProvider userProvider
}

func New(userProvider userProvider) *UseCase {
	return &UseCase{
		userProvider: userProvider,
	}
}
