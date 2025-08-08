package djstand

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

func (s *Stand) Skip(ctx context.Context, externalGroupID id.GroupExternal) error {
	err := s.skipUseCase.SkipByExternalGroupID(ctx, externalGroupID)
	if err != nil {
		return fmt.Errorf("failed to skip track: %w", err)
	}

	return nil
}
