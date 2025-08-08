package loopuc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/domain/entity/track"
)

type dj interface {
	NotifyQueueOrderType(ctx context.Context, externalGroupID id.GroupExternal, queueOrderType queue.OrderType) error
}

type queueProvider interface {
	Update(ctx context.Context, newQueue *queue.Queue) error
	GetByGroupID(ctx context.Context, groupID id.Group) (*queue.Queue, error)
}

type trackProvider interface {
	Get(ctx context.Context, trackID id.Track) (*track.Track, error)
}

type UseCase struct {
	dj            dj
	queueProvider queueProvider
	trackProvider trackProvider
}

func New(dj dj, queueProvider queueProvider, trackProvider trackProvider) *UseCase {
	return &UseCase{
		dj:            dj,
		queueProvider: queueProvider,
		trackProvider: trackProvider,
	}
}
