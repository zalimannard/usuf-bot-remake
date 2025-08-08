package dj

import (
	"fmt"
	"usuf-bot-remake/internal/domain/entity/track"
	"usuf-bot-remake/internal/domain/entity/user"
)

func (d *DJ) notifyTrackStarted(user *user.User, startedTrack *track.Track) error {
	url := startedTrack.URL()
	fmt.Println("Track started:", startedTrack.Title(), url.String(), user.Name())
	return nil
}
