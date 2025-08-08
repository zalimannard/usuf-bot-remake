package randomc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (c *Command) Execute(ctx context.Context, args []string) {
	err := c.randomOrderUseCase.Random(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(fmt.Errorf("failed to randomize track: %w", err)).Msg("Error")
	}
}
