package userrepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/user"
)

func (r *Repository) Create(ctx context.Context, userToCreate *user.User) error {
	r.userByID[userToCreate.ID()] = userToCreate
	r.userByExternalID[userToCreate.ExternalID()] = userToCreate

	return nil
}
