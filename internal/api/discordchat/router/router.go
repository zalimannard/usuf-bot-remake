package router

import (
	"context"
	"fmt"
	"os"
	"strings"
	"usuf-bot-remake/internal/api/discordchat/command"
	"usuf-bot-remake/internal/api/discordchat/middleware"
	"usuf-bot-remake/pkg/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Config interface {
	Prefix() string
}

type Router struct {
	prefix        string
	middleware    *middleware.Middleware
	executeByName map[string]func(ctx context.Context, args []string)
	log           *logger.Logger
}

func New(cfg Config, middleware *middleware.Middleware, commands []command.Command, log *logger.Logger) *Router {
	executeByName := make(map[string]func(ctx context.Context, args []string))
	for i := range commands {
		for j := range commands[i].Names() {
			executeByName[commands[i].Names()[j]] = commands[i].Execute
		}
	}

	return &Router{
		prefix:        cfg.Prefix(),
		middleware:    middleware,
		executeByName: executeByName,
		log:           log,
	}
}

func (r *Router) OnNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, r.prefix) || m.Author.Bot {
		return
	}

	parts := strings.Fields(m.Content[len(r.prefix):])
	if len(parts) == 0 {
		return
	}
	commandName := parts[0]
	args := parts[1:]

	if commandName == "reset" {
		fmt.Println("Resetting...")
		os.Exit(0)
	}

	executeCommand, isRegistered := r.executeByName[commandName]
	if !isRegistered {
		return
	}

	ctx := context.Background()
	ctx = r.log.WithContext(ctx)
	ctxWithInfo, err := r.middleware.RequesterInfo(ctx, m)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to attach requester info")
	}

	executeCommand(ctxWithInfo, args)
}
