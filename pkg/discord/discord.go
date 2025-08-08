package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrFailedToCreateDiscordSession = errors.New("failed to create discord session")
)

type Config interface {
	Token() string
}

type Discord struct {
	*discordgo.Session
}

func New(cfg Config) (*Discord, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.Token()))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToCreateDiscordSession, err.Error())
	}

	return &Discord{
		Session: session,
	}, nil
}
