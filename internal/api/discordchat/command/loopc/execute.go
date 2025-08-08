package loopc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (c *Command) Execute(ctx context.Context, args []string) {
	err := c.loopUseCase.Loop(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(fmt.Errorf("failed to loop track: %w", err)).Msg("Error")
	}
}
