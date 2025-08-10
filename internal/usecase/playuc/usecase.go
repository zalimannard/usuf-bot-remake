package playuc

import (
	"context"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/domain/entity/user"
)

type queueProvider interface {
	Update(ctx context.Context, newQueue *queue.Queue) error
	GetByGroupID(ctx context.Context, groupID id.Group) (*queue.Queue, error)
}

type trackProvider interface {
	GetByURL(ctx context.Context, targetURL url.URL) (*track.Track, error)
	ExpandURL(ctx context.Context, trackURL url.URL) ([]url.URL, error)
	GetURLByQuery(ctx context.Context, query string) (*url.URL, error)
}

type dj interface {
	Start(ctx context.Context, targetGroup *group.Group, targetUser *user.User, trackToStart *track.Track) error
}

type UseCase struct {
	queueProvider queueProvider
	trackProvider trackProvider
	dj            dj
}

func New(queueProvider queueProvider, trackProvider trackProvider, dj dj) *UseCase {
	return &UseCase{
		queueProvider: queueProvider,
		trackProvider: trackProvider,
		dj:            dj,
	}
}
