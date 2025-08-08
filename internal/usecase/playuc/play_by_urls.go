package playuc

import (
	"context"
	"fmt"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/util"
)

func (u *UseCase) PlayByURLs(ctx context.Context, trackURLs []url.URL) error {
	if len(trackURLs) == 0 {
		return nil
	}

	ctxGroup := ctx.Value(util.ContextKeyRequesterGroup).(*group.Group)
	ctxUser := ctx.Value(util.ContextKeyRequesterUser).(*user.User)

	currentQueue, err := u.queueProvider.GetByGroupID(ctx, ctxGroup.ID())
	if err != nil {
		return fmt.Errorf("failed to get queue by group id: %s", err.Error())
	}

	tracksToAdd := make([]track.Track, 0, len(trackURLs))
	for _, trackURL := range trackURLs {
		trackToAdd, err := u.trackProvider.GetByURL(ctx, trackURL)
		if err != nil {
			return fmt.Errorf("failed to get track by url: %s", err.Error())
		}
		tracksToAdd = append(tracksToAdd, *trackToAdd)
	}

	newQueueItems := make([]queue.Item, 0, len(tracksToAdd))
	for i := range tracksToAdd {
		newQueueItems = append(newQueueItems, queue.NewItem(nil, tracksToAdd[i].ID(), ctxUser.ID()))
	}

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
		err = u.dj.Start(ctx, ctxGroup, ctxUser, &tracksToAdd[0])
		if err != nil {
			return fmt.Errorf("failed to start track by dj: %s", err.Error())
		}
	}

	return nil
}
