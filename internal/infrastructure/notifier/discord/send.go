package discordnotifier

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/notification"

	"github.com/bwmarrin/discordgo"
)

func (n *Notifier) Send(ctx context.Context, channelID string, notificationToSend notification.Notification) error {
	title := ""
	if notificationToSend.Title != nil {
		title = *notificationToSend.Title
	}

	description := ""
	if notificationToSend.Description != nil {
		description = *notificationToSend.Description
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       3066993,
	}

	message := &discordgo.MessageSend{
		Embed: embed,
	}

	_, err := n.session.ChannelMessageSendComplex(channelID, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
