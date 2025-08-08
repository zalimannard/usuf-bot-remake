package useruc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/user"
)

func (u *UseCase) Update(ctx context.Context, userToUpdate *user.User) error {
	err := u.userProvider.Update(ctx, userToUpdate)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
