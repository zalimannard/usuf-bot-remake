package userprovider

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/infrastructure/repository/userrepo"
)

func (p *Provider) Update(ctx context.Context, userToUpdate *user.User) error {
	err := p.userRepository.Update(ctx, userToUpdate)
	if err != nil {
		if errors.Is(err, userrepo.ErrUserByExternalIDNotFound) {
			return newErrUserNotFound(userToUpdate.ID())
		}
		return fmt.Errorf("failed to update user in repository: %s", err.Error())
	}

	return nil
}
