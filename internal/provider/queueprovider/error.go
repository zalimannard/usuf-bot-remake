package queueprovider

import (
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrQueueByGroupIDNotFound = errors.New("queue by group id not found")
)

func newErrQueueByGroupIDNotFound(groupID id.Group) error {
	return fmt.Errorf("%w (group id=%s", ErrQueueByGroupIDNotFound, groupID.String())
}
