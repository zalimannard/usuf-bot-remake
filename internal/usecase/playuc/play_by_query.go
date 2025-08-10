package playuc

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/util"
)

func (u *UseCase) PlayByQuery(ctx context.Context, query string) error {
	if len(query) == 0 {
		return nil
	}

	ctxGroup := ctx.Value(util.ContextKeyRequesterGroup).(*group.Group)
	ctxUser := ctx.Value(util.ContextKeyRequesterUser).(*user.User)

	currentQueue, err := u.queueProvider.GetByGroupID(ctx, ctxGroup.ID())
	if err != nil {
		return fmt.Errorf("failed to get queue by group id: %s", err.Error())
	}

	urlToAdd, err := u.trackProvider.GetURLByQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to resolve tracks by query: %s", err.Error())
	}

	trackToAdd, err := u.trackProvider.GetByURL(ctx, *urlToAdd)
	if err != nil {
		return fmt.Errorf("failed to get track by url: %s", err.Error())
	}

	newQueueItems := make([]queue.Item, 0)
	newQueueItems = append(newQueueItems,
		queue.NewItem(nil, trackToAdd.ID(), ctxUser.ID()))

	currentNumber := currentQueue.CurrentNumber()
	if currentNumber == 0 {
		currentNumber = 1
	}
	newQueue, err := queue.New(
		util.Ptr(currentQueue.ID()),
		append(currentQueue.Items(), newQueueItems...),
		currentQueue.OrderType(),
		currentNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to construct new queue: %s", err.Error())
	}

	err = u.queueProvider.Update(ctx, newQueue)
	if err != nil {
		return fmt.Errorf("failed to update queue: %s", err.Error())
	}

	if currentQueue.CurrentNumber() == 0 {
		err = u.dj.Start(ctx, ctxGroup, ctxUser, trackToAdd)
		if err != nil {
			return fmt.Errorf("failed to start track by dj: %s", err.Error())
		}
	}

	return nil
}
