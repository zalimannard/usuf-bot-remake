package grouprepoinmemory

import (
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
)

type Repository struct {
	groupByID         map[id.Group]*group.Group
	groupByExternalID map[id.GroupExternal]*group.Group
}

func New() *Repository {
	return &Repository{
		groupByID:         make(map[id.Group]*group.Group),
		groupByExternalID: make(map[id.GroupExternal]*group.Group),
	}
}
