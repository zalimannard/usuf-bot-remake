package djstand

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
)

type skipUseCase interface {
	SkipByExternalGroupID(ctx context.Context, externalGroupID id.GroupExternal) error
}

type Stand struct {
	skipUseCase skipUseCase
}

func New(skipUseCase skipUseCase) *Stand {
	return &Stand{
		skipUseCase: skipUseCase,
	}
}

func (s *Stand) SetSkipUseCase(skipUseCase skipUseCase) {
	s.skipUseCase = skipUseCase
}
