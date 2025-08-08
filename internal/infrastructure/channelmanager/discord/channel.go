package discordchannelmanager

import (
	"context"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/infrastructure/channelmanager"
)

type Manager struct {
	channelIDByExternalGroupID map[id.GroupExternal]string
}

func New() *Manager {
	return &Manager{
		channelIDByExternalGroupID: make(map[id.GroupExternal]string),
	}
}

func (m *Manager) Set(ctx context.Context, externalGroupID id.GroupExternal, channelID string) {
	m.channelIDByExternalGroupID[externalGroupID] = channelID
}

func (m *Manager) Get(ctx context.Context, externalGroupID id.GroupExternal) (string, error) {
	channelID, ok := m.channelIDByExternalGroupID[externalGroupID]
	if !ok {
		return "", channelmanager.ErrNotFound
	}

	return channelID, nil
}
