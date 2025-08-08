package middleware

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
)

type groupUseCase interface {
	Create(ctx context.Context, groupToCreate *group.Group) error
	GetByExternalID(ctx context.Context, externalGroupID id.GroupExternal) (*group.Group, error)
}

type userUseCase interface {
	Create(ctx context.Context, userToCreate *user.User) error
	Update(ctx context.Context, userToUpdate *user.User) error
	GetByExternalID(ctx context.Context, externalUserID id.UserExternal) (*user.User, error)
}

type Middleware struct {
	groupUseCase groupUseCase
	userUseCase  userUseCase
}

func New(groupUseCase groupUseCase, userProvider userUseCase) *Middleware {
	return &Middleware{
		groupUseCase: groupUseCase,
		userUseCase:  userProvider,
	}
}
