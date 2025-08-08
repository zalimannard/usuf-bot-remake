package userprovider

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/infrastructure/repository/userrepo"
)

func (p *Provider) GetByExternalID(ctx context.Context, externalUserID id.UserExternal) (*user.User, error) {
	targetUser, err := p.userRepository.GetByExternalID(ctx, externalUserID)
	if err != nil {
		if errors.Is(err, userrepo.ErrUserByExternalIDNotFound) {
			return nil, newErrUserByExternalIDNotFound(externalUserID)
		}
		return nil, fmt.Errorf("failed to get user by external id in repository: %s", err.Error())
	}

	return targetUser, nil
}
