package grouprepo

import (
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrGroupNotFound             = errors.New("group not found in repository")
	ErrGroupByExternalIDNotFound = errors.New("group by external id not found in repository")
)

func NewErrGroupNotFound(groupID id.Group) error {
	return fmt.Errorf("%w (id=%s)", ErrGroupNotFound, groupID)
}

func NewErrGroupByExternalIDNotFound(externalGroupID id.GroupExternal) error {
	return fmt.Errorf("%w (external id=%s)", ErrGroupByExternalIDNotFound, externalGroupID)
}
