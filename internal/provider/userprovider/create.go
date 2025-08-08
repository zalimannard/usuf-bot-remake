package userprovider

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/user"
)

func (p *Provider) Create(ctx context.Context, userToCreate *user.User) error {
	err := p.userRepository.Create(ctx, userToCreate)
	if err != nil {
		return fmt.Errorf("failed to create user in repository: %s", err.Error())
	}

	return nil
}
