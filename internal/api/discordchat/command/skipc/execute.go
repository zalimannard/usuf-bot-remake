package skipc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (c *Command) Execute(ctx context.Context, args []string) {
	err := c.skipUseCase.Skip(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(fmt.Errorf("%w: %s", ErrFailedToSkipTrack, err.Error()))
	}
}
