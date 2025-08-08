package discordchat

import (
	"errors"
	"fmt"
	"usuf-bot-remake/pkg/discord"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrFailedToOpenDiscordSession  = errors.New("failed to open discord session")
	ErrFailedToCloseDiscordSession = errors.New("failed to close discord session")
)

type Router interface {
	OnNewMessage(s *discordgo.Session, m *discordgo.MessageCreate)
}

type Chat struct {
	session *discordgo.Session
	router  Router
}

func New(discord *discord.Discord, router Router) (*Chat, error) {

	return &Chat{
		session: discord.Session,
		router:  router,
	}, nil
}

func (c *Chat) Start() error {
	c.session.AddHandler(c.router.OnNewMessage)

	err := c.session.Open()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToOpenDiscordSession, err)
	}

	return nil
}

func (c *Chat) Stop() error {
	err := c.session.Close()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCloseDiscordSession, err)
	}

	return nil
}
