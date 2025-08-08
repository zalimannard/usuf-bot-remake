package clearc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (c *Command) Execute(ctx context.Context, args []string) {
	err := c.clearUseCase.Clear(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(fmt.Errorf("failed to clear track: %w", err)).Msg("Error")
	}
}
