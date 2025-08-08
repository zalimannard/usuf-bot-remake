package dj

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/helprow"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/notification"
	"usuf-bot-remake/internal/domain/entity/queue"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/infrastructure/channelmanager"
	"usuf-bot-remake/internal/util"
)

const (
	greenColor = 3066993
)

func (d *DJ) notifyTrackStarted(user *user.User, startedTrack *track.Track) error {
	url := startedTrack.URL()
	fmt.Println("Track started:", startedTrack.Title(), url.String(), user.Name())
	return nil
}

func (d *DJ) NotifyQueueOrderType(ctx context.Context, externalGroupID id.GroupExternal, queueOrderType queue.OrderType) error {
	channelID, err := d.channelManager.Get(ctx, externalGroupID)
	if err != nil {
		if errors.Is(err, channelmanager.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed getting channel id from external group: %w", err)
	}

	var orderName string
	switch queueOrderType {
	case queue.OrderTypeNormal:
		orderName = "По очереди"
	case queue.OrderTypeLoopTrack:
		orderName = "Трек зациклен"
	case queue.OrderTypeLoopQueue:
		orderName = "Очередь зациклена"
	case queue.OrderTypeRandom:
		orderName = "Случайный"
	default:
		return errors.New("unknown queue order type")
	}

	err = d.notifier.Send(ctx, channelID, []notification.Notification{
		{
			Title: util.Ptr(fmt.Sprintf("Порядок: %s", orderName)),
			Color: util.Ptr(greenColor),
		},
	})
	if err != nil {
		return fmt.Errorf("failed notifying channel: %w", err)
	}

	return nil
}

func (d *DJ) NotifyHelp(ctx context.Context, externalGroupID id.GroupExternal, rows []helprow.Row) error {
	channelID, err := d.channelManager.Get(ctx, externalGroupID)
	if err != nil {
		if errors.Is(err, channelmanager.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed getting channel id from external group: %w", err)
	}

	notifications := make([]notification.Notification, 0, len(rows))
	for _, row := range rows {
		notifications = append(notifications, notification.Notification{
			Title:       util.Ptr(row.Title),
			Description: util.Ptr(row.Description),
			Color:       util.Ptr(greenColor),
		})
	}

	err = d.notifier.Send(ctx, channelID, notifications)
	if err != nil {
		return fmt.Errorf("failed notifying channel: %w", err)
	}

	return nil
}
