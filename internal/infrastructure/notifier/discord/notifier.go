package discordnotifier

import "github.com/bwmarrin/discordgo"

type Notifier struct {
	session *discordgo.Session
}

func New(session *discordgo.Session) *Notifier {
	return &Notifier{
		session: session,
	}
}
