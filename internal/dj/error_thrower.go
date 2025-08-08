package dj

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/infrastructure/dancefloor"
)

func (d *DJ) ThrowError() {
	for {
		disorderEl := <-d.disorderChan
		if errors.Is(disorderEl.err, dancefloor.ErrEndOfTrack) {
			err := d.djStand.Skip(context.Background(), disorderEl.danceFloor.ExternalGroupID())
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}

		err := disorderEl.danceFloor.Abort()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
