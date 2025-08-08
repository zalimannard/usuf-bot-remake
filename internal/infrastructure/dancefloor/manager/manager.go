package dancefloormanager

import "github.com/bwmarrin/discordgo"

type Manager struct {
	session *discordgo.Session
}

func NewDiscord(session *discordgo.Session) *Manager {
	return &Manager{
		session: session,
	}
}
