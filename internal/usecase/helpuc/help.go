package helpuc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/helprow"
	"usuf-bot-remake/internal/util"
)

func (u *UseCase) Help(ctx context.Context, helpRows []helprow.Row) error {
	ctxGroup := ctx.Value(util.ContextKeyRequesterGroup).(*group.Group)

	return u.dj.NotifyHelp(ctx, ctxGroup.ExternalID(), helpRows)
}
