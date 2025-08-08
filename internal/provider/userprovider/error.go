package userprovider

import (
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserByExternalIDNotFound = errors.New("user by external id not found")
)

func newErrUserNotFound(userID id.User) error {
	return fmt.Errorf("%w (id=%s", ErrUserNotFound, userID.String())
}

func newErrUserByExternalIDNotFound(externalUserID id.UserExternal) error {
	return fmt.Errorf("%w (external id=%s", ErrUserByExternalIDNotFound, externalUserID.String())
}
