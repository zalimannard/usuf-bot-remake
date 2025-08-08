package userrepo

import (
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrUserNotFound                   = errors.New("user not found in repository")
	ErrUserByExternalIDNotFound       = errors.New("user by external id not found in repository")
	ErrUserIDAndExternalIDAreConflict = errors.New("user id and external id are conflict")
)

func NewErrUserNotFound(userID id.User) error {
	return fmt.Errorf("%w (id=%s)", ErrUserNotFound, userID)
}

func NewErrUserByExternalIDNotFound(externalUserID id.UserExternal) error {
	return fmt.Errorf("%w (external id=%s)", ErrUserByExternalIDNotFound, externalUserID)
}
