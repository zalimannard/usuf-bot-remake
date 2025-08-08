package id

import (
	"fmt"

	"github.com/google/uuid"
)

type Queue uuid.UUID

func ParseQueue(id string) (Queue, error) {
	queueID, err := uuid.Parse(id)
	if err != nil {
		return Queue(queueID), fmt.Errorf("invalid queue id: %w", err)
	}
	return Queue(queueID), nil
}

func GenerateQueue() Queue {
	return Queue(uuid.New())
}

func (g Queue) String() string {
	return uuid.UUID(g).String()
}

type QueueItem uuid.UUID

func ParseQueueItem(id string) (QueueItem, error) {
	queueItemID, err := uuid.Parse(id)
	if err != nil {
		return QueueItem(queueItemID), fmt.Errorf("invalid queue item id: %w", err)
	}
	return QueueItem(queueItemID), nil
}

func GenerateQueueItem() QueueItem {
	return QueueItem(uuid.New())
}

func (g QueueItem) String() string {
	return uuid.UUID(g).String()
}
