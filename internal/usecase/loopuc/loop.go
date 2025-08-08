package loopuc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/util"
)

func (u *UseCase) Loop(ctx context.Context) error {
	ctxGroup := ctx.Value(util.ContextKeyRequesterGroup).(*group.Group)

	currentQueue, err := u.queueProvider.GetByGroupID(ctx, ctxGroup.ID())
	if err != nil {
		return fmt.Errorf("failed to get queue by group id: %s", err.Error())
	}

	newOrderType := queue.OrderTypeLoopTrack
	if currentQueue.OrderType() == queue.OrderTypeLoopTrack {
		newOrderType = queue.OrderTypeNormal
	}

	newQueue, err := queue.New(
		util.Ptr(currentQueue.ID()),
		currentQueue.Items(),
		newOrderType,
		currentQueue.CurrentNumber(),
	)

	err = u.queueProvider.Update(ctx, newQueue)
	if err != nil {
		return fmt.Errorf("failed to update queue: %s", err.Error())
	}

	err = u.dj.NotifyQueueOrderType(ctx, ctxGroup.ExternalID(), newOrderType)
	if err != nil {
		return fmt.Errorf("failed to notify: %s", err.Error())
	}

	return nil
}
