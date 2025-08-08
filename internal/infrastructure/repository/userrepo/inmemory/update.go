package userrepoinmemory

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/infrastructure/repository/userrepo"
)

func (r *Repository) Update(ctx context.Context, userToUpdate *user.User) error {
	currentUserByID, existsByID := r.userByID[userToUpdate.ID()]
	if !existsByID {
		return userrepo.NewErrUserNotFound(userToUpdate.ID())
	}
	currentUserByExternalID, existsByExternalID := r.userByExternalID[userToUpdate.ExternalID()]
	if !existsByExternalID {
		return userrepo.NewErrUserByExternalIDNotFound(userToUpdate.ExternalID())
	}
	if currentUserByID.ID() != currentUserByExternalID.ID() {
		return userrepo.ErrUserIDAndExternalIDAreConflict
	}

	r.userByID[userToUpdate.ID()] = userToUpdate
	r.userByExternalID[userToUpdate.ExternalID()] = userToUpdate

	return nil
}
