package userprovider

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/infrastructure/repository/userrepo"
)

func (p *Provider) Get(ctx context.Context, userID id.User) (*user.User, error) {
	targetUser, err := p.userRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, userrepo.ErrUserNotFound) {
			return nil, newErrUserNotFound(userID)
		}
		return nil, fmt.Errorf("failed to get user in repository: %s", err.Error())
	}

	return targetUser, nil
}
