package clearuc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/util"
)

func (u *UseCase) Clear(ctx context.Context) error {
	ctxGroup := ctx.Value(util.ContextKeyRequesterGroup).(*group.Group)

	currentQueue, err := u.queueProvider.GetByGroupID(ctx, ctxGroup.ID())
	if err != nil {
		return fmt.Errorf("failed to get queue by group id: %s", err.Error())
	}

	err = u.queueProvider.Delete(ctx, currentQueue.ID())
	if err != nil {
		return fmt.Errorf("failed to update queue: %s", err.Error())
	}

	err = u.dj.Close(ctx, ctxGroup.ID())
	if err != nil {
		return fmt.Errorf("failed to stop by dj: %s", err.Error())
	}

	err = u.dj.NotifyClearQueue(ctx, ctxGroup.ExternalID())
	if err != nil {
		return fmt.Errorf("failed to notify: %s", err.Error())
	}

	return nil
}
