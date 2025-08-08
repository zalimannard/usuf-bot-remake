package loopqc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (c *Command) Execute(ctx context.Context, args []string) {
	err := c.loopqUseCase.Loopq(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(fmt.Errorf("failed to loop queue: %w", err)).Msg("Error")
	}
}
