package clearuc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
)

type dj interface {
	NotifyClearQueue(ctx context.Context, externalGroupID id.GroupExternal) error
	Close(ctx context.Context, groupID id.Group) error
}

type queueProvider interface {
	GetByGroupID(ctx context.Context, groupID id.Group) (*queue.Queue, error)
	Delete(ctx context.Context, queueID id.Queue) error
}

type UseCase struct {
	dj            dj
	queueProvider queueProvider
}

func New(dj dj, queueProvider queueProvider) *UseCase {
	return &UseCase{
		dj:            dj,
		queueProvider: queueProvider,
	}
}
