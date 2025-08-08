package dj

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/id"
)

func (d *DJ) Close(ctx context.Context, groupID id.Group) error {
	targetDanceFloor, exists := d.danceFloorByGroupID[groupID]
	if !exists {
		return fmt.Errorf("dance floor not found")
	}

	delete(d.danceFloorByGroupID, groupID)

	err := targetDanceFloor.Abort()
	if err != nil {
		return fmt.Errorf("failed to stop tracks: %s", err.Error())
	}

	err = targetDanceFloor.Close()
	if err != nil {
		return fmt.Errorf("failed to close dance floor: %s", err.Error())
	}

	return nil
}
