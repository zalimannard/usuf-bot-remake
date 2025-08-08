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
	"usuf-bot-remake/internal/api/discordchat/command/clearc"
	"usuf-bot-remake/internal/api/discordchat/command/helpc"
	"usuf-bot-remake/internal/api/discordchat/command/loopc"
	"usuf-bot-remake/internal/api/discordchat/command/loopqc"
	"usuf-bot-remake/internal/api/discordchat/command/playc"
	"usuf-bot-remake/internal/api/discordchat/command/randomc"
	"usuf-bot-remake/internal/api/discordchat/command/skipc"
	"usuf-bot-remake/internal/api/discordchat/middleware"
	"usuf-bot-remake/internal/api/discordchat/router"
	"usuf-bot-remake/internal/app"
	discordchannelmanager "usuf-bot-remake/internal/infrastructure/channelmanager/discord"
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

	channelManager := discordchannelmanager.New()

	application := app.New(discordSession, channelManager)

	playCommand := playc.New(application.PlayUseCase())
	skipCommand := skipc.New(application.SkipUseCase())
	loopCommand := loopc.New(application.LoopUseCase())
	loopqCommand := loopqc.New(application.LoopqUseCase())
	randomCommand := randomc.New(application.RandomUseCase())
	clearCommand := clearc.New(application.ClearUseCase())
	helpCommand := helpc.New(application.HelpUseCase())

	discordMiddleware := middleware.New(
		application.GroupUseCase(),
		application.UserUseCase(),
	)

	discordChatRouter := router.New(cfg.Discord(), channelManager, discordMiddleware,
		[]command.Command{
			playCommand,
			skipCommand,
			loopCommand,
			loopqCommand,
			randomCommand,
			clearCommand,
			helpCommand,
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
