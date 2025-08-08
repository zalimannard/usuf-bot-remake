package queuerepoinmemory

import (
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
)

type Repository struct {
	queueByID        map[id.Queue]*queue.Queue
	queueByGroupID   map[id.Group]*queue.Queue
	groupIDByQueueID map[id.Queue]id.Group
	queueIDByGroupID map[id.Group]id.Queue
}

func New() *Repository {
	return &Repository{
		queueByGroupID:   make(map[id.Group]*queue.Queue),
		groupIDByQueueID: make(map[id.Queue]id.Group),
		queueIDByGroupID: make(map[id.Group]id.Queue),
		queueByID:        make(map[id.Queue]*queue.Queue),
	}
}
