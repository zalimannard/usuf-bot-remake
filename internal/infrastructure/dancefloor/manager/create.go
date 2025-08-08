package dancefloormanager

import (
	"fmt"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/interface/dancefloor"
	dancefloordiscord "usuf-bot-remake/internal/infrastructure/dancefloor/discord"
)

type DanceFloor interface {
	Play(urlToPlay url.URL) error
	Abort() error
	ErrChan() <-chan error
}

func (m *Manager) Create(groupExternalID id.GroupExternal, userExternalID id.UserExternal) (dancefloor.DanceFloor, error) {
	guild, err := m.session.State.Guild(groupExternalID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get guild: %w", err)
	}
	var channelID string
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userExternalID.String() {
			channelID = vs.ChannelID
		}
	}

	return dancefloordiscord.New(m.session, groupExternalID.String(), channelID), nil
}
