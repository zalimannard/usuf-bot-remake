package app

import (
	"usuf-bot-remake/internal/api/djstand"
	"usuf-bot-remake/internal/dj"
	dancefloormanager "usuf-bot-remake/internal/infrastructure/dancefloor/manager"
	discordnotifier "usuf-bot-remake/internal/infrastructure/notifier/discord"
	grouprepoinmemory "usuf-bot-remake/internal/infrastructure/repository/grouprepo/inmemory"
	queuerepoinmemory "usuf-bot-remake/internal/infrastructure/repository/queuerepo/inmemory"
	trackrepoinmemory "usuf-bot-remake/internal/infrastructure/repository/trackrepo/inmemory"
	userrepoinmemory "usuf-bot-remake/internal/infrastructure/repository/userrepo/inmemory"
	"usuf-bot-remake/internal/infrastructure/trackloader"
	"usuf-bot-remake/internal/provider/groupprovider"
	"usuf-bot-remake/internal/provider/queueprovider"
	"usuf-bot-remake/internal/provider/trackprovider"
	"usuf-bot-remake/internal/provider/userprovider"
	"usuf-bot-remake/internal/usecase/groupuc"
	"usuf-bot-remake/internal/usecase/playuc"
	"usuf-bot-remake/internal/usecase/skipuc"
	"usuf-bot-remake/internal/usecase/useruc"
	"usuf-bot-remake/pkg/discord"
)

type Application struct {
	groupUseCase *groupuc.UseCase
	userUseCase  *useruc.UseCase
	playUseCase  *playuc.UseCase
	skipUseCase  *skipuc.UseCase
}

func New(session *discord.Discord) *Application {
	groupRepository := grouprepoinmemory.New()
	groupProvider := groupprovider.New(groupRepository)

	queueRepository := queuerepoinmemory.New()
	queueProvider := queueprovider.New(queueRepository)

	trackLoader := trackloader.New()
	trackRepository := trackrepoinmemory.New()
	trackProvider := trackprovider.New(trackRepository, trackLoader)

	userRepository := userrepoinmemory.New()
	userProvider := userprovider.New(userRepository)

	djStand := djstand.New(nil)

	danceFloorManager := dancefloormanager.NewDiscord(session.Session)
	notifier := discordnotifier.New()
	diskJockey := dj.New(djStand, danceFloorManager, notifier)

	skipUseCase := skipuc.New(groupProvider, diskJockey, queueProvider, trackProvider, userProvider)
	groupUseCase := groupuc.New(groupProvider)
	userUseCase := useruc.New(userProvider)
	playUseCase := playuc.New(queueProvider, trackProvider, diskJockey)

	djStand.SetSkipUseCase(skipUseCase)

	return &Application{
		groupUseCase: groupUseCase,
		userUseCase:  userUseCase,
		playUseCase:  playUseCase,
		skipUseCase:  skipUseCase,
	}
}

func (a *Application) GroupUseCase() *groupuc.UseCase {
	return a.groupUseCase
}

func (a *Application) UserUseCase() *useruc.UseCase {
	return a.userUseCase
}

func (a *Application) PlayUseCase() *playuc.UseCase {
	return a.playUseCase
}

func (a *Application) SkipUseCase() *skipuc.UseCase {
	return a.skipUseCase
}
