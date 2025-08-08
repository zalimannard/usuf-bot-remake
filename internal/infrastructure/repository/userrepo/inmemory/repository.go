package userrepoinmemory

import (
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
)

type Repository struct {
	userByID         map[id.User]*user.User
	userByExternalID map[id.UserExternal]*user.User
}

func New() *Repository {
	return &Repository{
		userByID:         make(map[id.User]*user.User),
		userByExternalID: make(map[id.UserExternal]*user.User),
	}
}
