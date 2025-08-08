package user

import (
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/util"
)

type User struct {
	id         id.User
	externalID id.UserExternal
	name       string
}

func New(userID *id.User, externalID id.UserExternal, name string) *User {
	if userID == nil {
		userID = util.Ptr(id.GenerateUser())
	}
	return &User{
		id:         *userID,
		externalID: externalID,
		name:       name,
	}
}

func (u *User) ID() id.User {
	return u.id
}

func (u *User) ExternalID() id.UserExternal {
	return u.externalID
}

func (u *User) Name() string {
	return u.name
}
