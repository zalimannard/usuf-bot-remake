package playuc

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/util"
)

var groupLocks sync.Map

func groupLock(gid id.Group) *sync.Mutex {
	l, _ := groupLocks.LoadOrStore(gid, &sync.Mutex{})
	return l.(*sync.Mutex)
}

func (u *UseCase) PlayByURLs(ctx context.Context, trackURLs []url.URL) error {
	if len(trackURLs) == 0 {
		return nil
	}

	ctxGroup := ctx.Value(util.ContextKeyRequesterGroup).(*group.Group)
	ctxUser := ctx.Value(util.ContextKeyRequesterUser).(*user.User)

	// Сериализация операций для одной группы
	lock := groupLock(ctxGroup.ID())
	lock.Lock()
	defer lock.Unlock()

	currentQueue, err := u.queueProvider.GetByGroupID(ctx, ctxGroup.ID())
	if err != nil {
		return fmt.Errorf("failed to get queue by group id: %s", err.Error())
	}

	tracksToAdd := make([]track.Track, 0, len(trackURLs))
	for _, trackURL := range trackURLs {
		resolvedURL, err := u.trackProvider.ExpandURL(ctx, trackURL)
		if err != nil {
			return fmt.Errorf("failed to resolve tracks by url: %s", err.Error())
		}

		for i := range resolvedURL {
			trackToAdd, err := u.trackProvider.GetByURL(ctx, resolvedURL[i])
			if err != nil {
				return fmt.Errorf("failed to get track by url: %s", err.Error())
			}
			tracksToAdd = append(tracksToAdd, *trackToAdd)
		}
	}

	newQueueItems := make([]queue.Item, 0, len(tracksToAdd))
	for i := range tracksToAdd {
		newQueueItems = append(newQueueItems, queue.NewItem(nil, tracksToAdd[i].ID(), ctxUser.ID()))
	}

	// Запомним, была ли очередь пустой перед обновлением — только первый вызов внутри мьютекса увидит 0
	wasIdle := currentQueue.CurrentNumber() == 0

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

	if err = u.queueProvider.Update(ctx, newQueue); err != nil {
		return fmt.Errorf("failed to update queue: %s", err.Error())
	}

	// Стартуем воспроизведение только если очередь была пуста до этого обновления
	if wasIdle && len(tracksToAdd) > 0 {
		if err = u.dj.Start(ctx, ctxGroup, ctxUser, &tracksToAdd[0]); err != nil {
			return fmt.Errorf("failed to start track by dj: %s", err.Error())
		}
	}

	for i := range tracksToAdd {
		if err = u.dj.NotifyTrackAdded(ctx, ctxGroup.ExternalID(), currentQueue.Length()+1+i, newQueue.Length(), *ctxUser, tracksToAdd[i]); err != nil {
			return fmt.Errorf("failed to notify: %s", err.Error())
		}
	}

	if wasIdle && len(tracksToAdd) > 0 {
		if err = u.dj.NotifyNowPlaying(ctx, ctxGroup.ExternalID(), 1, newQueue.Length(), *ctxUser, tracksToAdd[0]); err != nil {
			return fmt.Errorf("failed to notify: %s", err.Error())
		}
	}

	return nil
}
