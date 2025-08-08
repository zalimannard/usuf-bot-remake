package dancefloordiscord

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/infrastructure/dancefloor"
	"usuf-bot-remake/pkg/streamer"

	"github.com/bwmarrin/discordgo"
)

type DanceFloor struct {
	session                 *discordgo.Session
	voiceConnection         *discordgo.VoiceConnection
	playBackgroundCtx       context.Context
	playBackgroundCtxCancel context.CancelFunc
	guildID                 string
	channelID               string
	errChan                 chan error
}

func New(session *discordgo.Session, guildID string, channelID string) *DanceFloor {
	return &DanceFloor{
		session:   session,
		guildID:   guildID,
		channelID: channelID,
		errChan:   make(chan error, 1),
	}
}

func (d *DanceFloor) ExternalGroupID() id.GroupExternal {
	return id.ParseGroupExternal(d.guildID)
}

func (d *DanceFloor) Play(urlToPlay url.URL) error {
	if d.voiceConnection == nil {
		newVoiceConnection, err := d.session.ChannelVoiceJoin(d.guildID, d.channelID, false, false)
		if err != nil {
			return fmt.Errorf("failed to join voice channel: %w", err)
		}
		d.voiceConnection = newVoiceConnection
	}

	err := d.voiceConnection.Speaking(true)
	if err != nil {
		return fmt.Errorf("failed to start speaking: %w", err)
	}

	d.playBackgroundCtx, d.playBackgroundCtxCancel = context.WithCancel(context.Background())

	go d.playBackground(d.playBackgroundCtx, urlToPlay)

	return nil
}

func (d *DanceFloor) playBackground(ctx context.Context, urlToPlay url.URL) {
	opusChan, errChan := streamer.Play(ctx, urlToPlay)
	for {
		select {
		case <-ctx.Done():
			return
		case contentPart, opened := <-opusChan:
			if !opened {
				return
			}
			d.voiceConnection.OpusSend <- contentPart
		case err := <-errChan:
			if errors.Is(err, streamer.ErrEndOfStream) {
				d.errChan <- fmt.Errorf("%w: %s", dancefloor.ErrEndOfTrack, err.Error())
			} else if errors.Is(err, streamer.ErrPlaybackUnavailable) {
				d.errChan <- fmt.Errorf("%w: %s", dancefloor.ErrEndOfTrack, err.Error())
			} else {
				d.errChan <- fmt.Errorf("error while playing: %w", err)
			}
			return
		}
	}
}

func (d *DanceFloor) Abort() error {
	if d.voiceConnection != nil {
		err := d.voiceConnection.Speaking(false)
		if err != nil {
			fmt.Println("failed to stop speaking: %w", err)
		}
	}

	if d.playBackgroundCtxCancel != nil {
		d.playBackgroundCtxCancel()
	}
	d.playBackgroundCtx = nil
	d.playBackgroundCtxCancel = nil

	return nil
}

func (d *DanceFloor) Close() error {
	err := d.voiceConnection.Disconnect()
	if err != nil {
		return fmt.Errorf("failed to disconnect voice connection: %w", err)
	}

	return nil
}

func (d *DanceFloor) ErrChan() <-chan error {
	return d.errChan
}
