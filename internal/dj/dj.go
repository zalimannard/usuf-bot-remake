package dj

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/interface/dancefloor"
)

type djStand interface {
	Skip(ctx context.Context, externalGroupID id.GroupExternal) error
}

type danceFloorManager interface {
	Create(groupID id.GroupExternal, userID id.UserExternal) (dancefloor.DanceFloor, error)
}

type notifier interface {
}

type DJ struct {
	djStand             djStand
	danceFloorManager   danceFloorManager
	notifier            notifier
	danceFloorByGroupID map[id.Group]dancefloor.DanceFloor
	disorderChan        chan disorder
}

func New(djStand djStand, danceFloorManager danceFloorManager, notifier notifier) *DJ {
	targetDJ := &DJ{
		djStand:             djStand,
		danceFloorManager:   danceFloorManager,
		notifier:            notifier,
		danceFloorByGroupID: make(map[id.Group]dancefloor.DanceFloor),
		disorderChan:        make(chan disorder, 1),
	}

	go targetDJ.ThrowError()

	return targetDJ
}

type disorder struct {
	danceFloor dancefloor.DanceFloor
	err        error
}
