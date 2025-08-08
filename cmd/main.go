package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"usuf-bot-remake/config"
	"usuf-bot-remake/internal/api/discordchat"
	"usuf-bot-remake/internal/api/discordchat/command"
	"usuf-bot-remake/internal/api/discordchat/command/playc"
	"usuf-bot-remake/internal/api/discordchat/command/skipc"
	"usuf-bot-remake/internal/api/discordchat/middleware"
	"usuf-bot-remake/internal/api/discordchat/router"
	"usuf-bot-remake/internal/app"
	"usuf-bot-remake/pkg/discord"
	"usuf-bot-remake/pkg/logger"

	"github.com/lrstanley/go-ytdlp"
)

func main() {
	ytdlp.MustInstall(context.TODO(), nil)

	cfg, err := config.Parse()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	log, err := logger.New(cfg.Logger())
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}

	discordSession, err := discord.New(cfg.Discord())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create discord session")
	}
	defer discordSession.Close()

	application := app.New(discordSession)

	playCommand := playc.New(application.PlayUseCase())
	skipCommand := skipc.New(application.SkipUseCase())

	discordMiddleware := middleware.New(
		application.GroupUseCase(),
		application.UserUseCase(),
	)

	discordChatRouter := router.New(cfg.Discord(), discordMiddleware,
		[]command.Command{
			playCommand,
			skipCommand,
		},
		log,
	)

	discordChat, err := discordchat.New(discordSession, discordChatRouter)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create discord chat")
	}
	err = discordChat.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start discord chat")
	}
	defer func(chat *discordchat.Chat) {
		err := chat.Stop()
		if err != nil {
			log.Error().Err(err).Msg("Failed to stop discord chat")
		}
	}(discordChat)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
