package groupprovider

import (
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrGroupNotFound             = errors.New("group not found")
	ErrGroupByExternalIDNotFound = errors.New("group by external id not found")
)

func newErrGroupNotFound(groupID id.Group) error {
	return fmt.Errorf("%w (id=%s", ErrGroupNotFound, groupID.String())
}

func newErrGroupByExternalIDNotFound(externalGroupID id.GroupExternal) error {
	return fmt.Errorf("%w (external id=%s", ErrGroupByExternalIDNotFound, externalGroupID.String())
}
