package userprovider

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
)

type Provider struct {
	userRepository userRepository
}

func New(userRepository userRepository) *Provider {
	return &Provider{
		userRepository: userRepository,
	}
}

type userRepository interface {
	Create(ctx context.Context, userToCreate *user.User) error
	Update(ctx context.Context, userToUpdate *user.User) error
	Get(ctx context.Context, userID id.User) (*user.User, error)
	GetByExternalID(ctx context.Context, externalUserID id.UserExternal) (*user.User, error)
}
