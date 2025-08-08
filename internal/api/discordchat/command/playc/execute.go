package playc

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

func (c *Command) Execute(ctx context.Context, args []string) {
	if len(args) == 0 {
		log.Ctx(ctx).Error().Msg("Invalid arguments")
		return
	}

	urls, argsIsURLs := isURLs(args)
	if argsIsURLs {
		err := c.playUseCase.PlayByURLs(ctx, urls)
		if err != nil {
			log.Ctx(ctx).Error().Err(fmt.Errorf("failed to play by urls: %w", err)).Msg("Error")
		}
	} else {
		err := c.playUseCase.PlayByQuery(ctx, strings.Join(args, " "))
		if err != nil {
			log.Ctx(ctx).Error().Err(fmt.Errorf("failed to play by query: %w", err)).Msg("Error")
		}
	}
}

func isURLs(args []string) ([]url.URL, bool) {
	if len(args) == 0 {
		return nil, false
	}

	urls := make([]url.URL, 0)

	for _, arg := range args {
		oneUrl, err := url.Parse(arg)
		if err != nil || oneUrl.Scheme == "" {
			return nil, false
		}

		urls = append(urls, *oneUrl)
	}

	return urls, true
}
