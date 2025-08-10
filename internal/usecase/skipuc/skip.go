package skipuc

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/util"
)

func (u *UseCase) Skip(ctx context.Context) error {
	ctxGroup := ctx.Value(util.ContextKeyRequesterGroup).(*group.Group)

	return u.skip(ctx, ctxGroup.ID())
}

func (u *UseCase) SkipByExternalGroupID(ctx context.Context, externalGroupID id.GroupExternal) error {
	targetGroup, err := u.groupProvider.GetByExternalID(ctx, externalGroupID)
	if err != nil {
		return fmt.Errorf("failed to get target group: %s", err.Error())
	}

	return u.skip(ctx, targetGroup.ID())
}

func (u *UseCase) skip(ctx context.Context, groupID id.Group) error {
	currentQueue, err := u.queueProvider.GetByGroupID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get queue by group id: %s", err.Error())
	}

	if currentQueue.OrderType() == queue.OrderTypeNormal && currentQueue.CurrentNumber() == currentQueue.Length() {
		err = u.queueProvider.Delete(ctx, currentQueue.ID())
		if err != nil {
			return fmt.Errorf("failed to delete queue")
		}

		err = u.dj.Close(ctx, groupID)
		if err != nil {
			return fmt.Errorf("failed to stop by dj: %s", err.Error())
		}

		return nil
	}

	var nextNumber int
	switch currentQueue.OrderType() {
	case queue.OrderTypeNormal:
		nextNumber = currentQueue.CurrentNumber() + 1
	case queue.OrderTypeLoopTrack:
		nextNumber = currentQueue.CurrentNumber()
	case queue.OrderTypeLoopQueue:
		if currentQueue.CurrentNumber() == currentQueue.Length() {
			nextNumber = 1
		} else {
			nextNumber = currentQueue.CurrentNumber() + 1
		}
	case queue.OrderTypeRandom:
		if currentQueue.Length() == 1 {
			nextNumber = 1
			break
		}
		generatedNumber := currentQueue.CurrentNumber()
		for generatedNumber == currentQueue.CurrentNumber() {
			rand.Seed(time.Now().UnixNano())
			generatedNumber = rand.Intn(currentQueue.Length()) + 1
		}
		nextNumber = generatedNumber
	default:
		return fmt.Errorf("unknown order type: %s", currentQueue.OrderType())
	}

	newQueue, err := queue.New(
		util.Ptr(currentQueue.ID()),
		currentQueue.Items(),
		currentQueue.OrderType(),
		nextNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to construct new queue: %s", err.Error())
	}

	nextItem := newQueue.Items()[nextNumber-1]
	nextTrack, err := u.trackProvider.Get(ctx, nextItem.TrackID())
	if err != nil {
		return fmt.Errorf("failed to get next track: %s", err.Error())
	}

	targetGroup, err := u.groupProvider.Get(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get target group: %s", err.Error())
	}
	trackRequester, err := u.userProvider.Get(ctx, nextItem.RequesterID())
	if err != nil {
		return fmt.Errorf("failed to get track requester: %s", err.Error())
	}

	requester, err := u.userProvider.Get(ctx, trackRequester.ID())
	if err != nil {
		return fmt.Errorf("failed to get requester: %s", err.Error())
	}

	err = u.queueProvider.Update(ctx, newQueue)
	if err != nil {
		return fmt.Errorf("failed to update queue: %s", err.Error())
	}

	err = u.dj.Start(ctx, targetGroup, trackRequester, nextTrack)
	if err != nil {
		return fmt.Errorf("failed to start new track by dj: %s", err.Error())
	}

	err = u.dj.NotifyNowPlaying(ctx, targetGroup.ExternalID(), nextNumber, newQueue.Length(), *requester, *nextTrack)
	if err != nil {
		return fmt.Errorf("failed to notify: %s", err.Error())
	}

	return nil
}
