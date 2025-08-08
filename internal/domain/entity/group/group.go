package group

import (
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/util"
)

type Group struct {
	id         id.Group
	externalID id.GroupExternal
}

func New(groupID *id.Group, externalID id.GroupExternal) *Group {
	if groupID == nil {
		groupID = util.Ptr(id.GenerateGroup())
	}
	return &Group{
		id:         *groupID,
		externalID: externalID,
	}
}

func (g *Group) ID() id.Group {
	return g.id
}

func (g *Group) ExternalID() id.GroupExternal {
	return g.externalID
}
