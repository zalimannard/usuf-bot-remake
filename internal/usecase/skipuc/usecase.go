package skipuc

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/domain/entity/user"
)

type groupProvider interface {
	Get(ctx context.Context, groupID id.Group) (*group.Group, error)
	GetByExternalID(ctx context.Context, externalGroupID id.GroupExternal) (*group.Group, error)
}

type dj interface {
	Start(ctx context.Context, targetGroup *group.Group, targetUser *user.User, trackToStart *track.Track) error
	Close(ctx context.Context, groupID id.Group) error
}

type queueProvider interface {
	Update(ctx context.Context, newQueue *queue.Queue) error
	GetByGroupID(ctx context.Context, groupID id.Group) (*queue.Queue, error)
	Delete(ctx context.Context, queueID id.Queue) error
}

type trackProvider interface {
	Get(ctx context.Context, trackID id.Track) (*track.Track, error)
}

type userProvider interface {
	Get(ctx context.Context, userID id.User) (*user.User, error)
}

type UseCase struct {
	groupProvider groupProvider
	dj            dj
	queueProvider queueProvider
	trackProvider trackProvider
	userProvider  userProvider
}

func New(groupProvider groupProvider, dj dj, queueProvider queueProvider, trackProvider trackProvider, userProvider userProvider) *UseCase {
	return &UseCase{
		groupProvider: groupProvider,
		dj:            dj,
		queueProvider: queueProvider,
		trackProvider: trackProvider,
		userProvider:  userProvider,
	}
}
