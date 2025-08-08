package discordnotifier

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/notification"

	"github.com/bwmarrin/discordgo"
)

func (n *Notifier) Send(ctx context.Context, channelID string, notificationsToSend []notification.Notification) error {
	if len(notificationsToSend) == 0 {
		return nil
	} else if len(notificationsToSend) == 1 {
		return n.sendOne(ctx, channelID, notificationsToSend[0])
	} else {
		return n.sendMany(ctx, channelID, notificationsToSend)
	}
}

func (n *Notifier) sendOne(ctx context.Context, channelID string, notificationToSend notification.Notification) error {
	title := ""
	if notificationToSend.Title != nil {
		title = *notificationToSend.Title
	}

	description := ""
	if notificationToSend.Description != nil {
		description = *notificationToSend.Description
	}

	color := 0
	if notificationToSend.Color != nil {
		color = *notificationToSend.Color
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
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

func (n *Notifier) sendMany(ctx context.Context, channelID string, notificationToSend []notification.Notification) error {
	color := 0
	if notificationToSend[0].Color != nil {
		color = *notificationToSend[0].Color
	}

	embeds := make([]*discordgo.MessageEmbedField, 0, len(notificationToSend))
	for _, notification := range notificationToSend {
		title := ""
		if notification.Title != nil {
			title = *notification.Title
		}

		description := ""
		if notification.Description != nil {
			description = *notification.Description
		}

		embeds = append(embeds, &discordgo.MessageEmbedField{
			Name:  title,
			Value: description,
		})
	}

	message := &discordgo.MessageEmbed{
		Fields: embeds,
		Color:  color,
	}

	_, err := n.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: message,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
