package id

import (
	"fmt"

	"github.com/google/uuid"
)

type Group uuid.UUID

func ParseGroup(id string) (Group, error) {
	groupID, err := uuid.Parse(id)
	if err != nil {
		return Group(groupID), fmt.Errorf("invalid group id: %w", err)
	}
	return Group(groupID), nil
}

func GenerateGroup() Group {
	return Group(uuid.New())
}

func (g Group) String() string {
	return uuid.UUID(g).String()
}

type GroupExternal string

func ParseGroupExternal(id string) GroupExternal {
	return GroupExternal(id)
}

func (g GroupExternal) String() string {
	return string(g)
}
