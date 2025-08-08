package helpuc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/helprow"
	"usuf-bot-remake/internal/domain/entity/id"
)

type dj interface {
	NotifyHelp(ctx context.Context, externalGroupID id.GroupExternal, rows []helprow.Row) error
}

type UseCase struct {
	dj dj
}

func New(dj dj) *UseCase {
	return &UseCase{
		dj: dj,
	}
}
