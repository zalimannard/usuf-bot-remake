package dancefloor

import (
	"net/url"
	"usuf-bot-remake/internal/domain/entity/id"
)

type DanceFloor interface {
	ExternalGroupID() id.GroupExternal
	Play(urlToPlay url.URL) error
	Abort() error
	Close() error
	ErrChan() <-chan error
}

type Manager interface {
	Create(groupID id.GroupExternal, userID id.UserExternal) (DanceFloor, error)
}
