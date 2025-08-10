package dj

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/domain/interface/dancefloor"
)

func (d *DJ) Start(ctx context.Context, targetGroup *group.Group, targetUser *user.User, trackToStart *track.Track) error {
	targetDanceFloor, exists := d.danceFloorByGroupID[targetGroup.ID()]
	if !exists {
		newDanceFloor, err := d.danceFloorManager.Create(targetGroup.ExternalID(), targetUser.ExternalID())
		if err != nil {
			return fmt.Errorf("failed to create dance floor: %s", err.Error())
		}
		go func(danceFloor dancefloor.DanceFloor) {
			// TODO: Сделать завершение
			for {
				err := <-newDanceFloor.ErrChan()
				d.disorderChan <- disorder{
					danceFloor: danceFloor,
					err:        err,
				}
			}
		}(newDanceFloor)
		targetDanceFloor = newDanceFloor
		d.danceFloorByGroupID[targetGroup.ID()] = targetDanceFloor
	}

	err := targetDanceFloor.Abort()
	if err != nil {
		return fmt.Errorf("failed to stop tracks: %s", err.Error())
	}

	err = targetDanceFloor.Play(trackToStart.URL())
	if err != nil {
		return fmt.Errorf("failed to play track: %s", err.Error())
	}

	err = d.notifyTrackStarted(targetUser, trackToStart)
	if err != nil {
		return fmt.Errorf("failed to notify track started: %s", err.Error())
	}

	return nil
}
