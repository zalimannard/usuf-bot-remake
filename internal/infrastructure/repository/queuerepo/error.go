package queuerepo

import (
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

var (
	ErrQueueNotFound          = errors.New("queue  not found")
	ErrQueueByGroupIDNotFound = errors.New("queue by group id not found")
)

func NewErrQueueNotFound(queueID id.Queue) error {
	return fmt.Errorf("%w (id=%s)", ErrQueueNotFound, queueID)
}

func NewErrQueueByGroupIDNotFound(groupID id.Group) error {
	return fmt.Errorf("%w (group id=%s)", ErrQueueByGroupIDNotFound, groupID)
}
