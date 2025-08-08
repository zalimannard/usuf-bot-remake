package id

import (
	"fmt"

	"github.com/google/uuid"
)

type User uuid.UUID

func ParseUser(id string) (User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return User(userID), fmt.Errorf("invalid user id: %w", err)
	}
	return User(userID), nil
}

func GenerateUser() User {
	return User(uuid.New())
}

func (g User) String() string {
	return uuid.UUID(g).String()
}

type UserExternal string

func ParseUserExternal(id string) UserExternal {
	return UserExternal(id)
}

func (g UserExternal) String() string {
	return string(g)
}
