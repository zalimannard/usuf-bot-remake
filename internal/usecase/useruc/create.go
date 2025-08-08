package useruc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/user"
)

func (u *UseCase) Create(ctx context.Context, userToCreate *user.User) error {
	err := u.userProvider.Create(ctx, userToCreate)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
