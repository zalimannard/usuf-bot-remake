package track

import (
	"net/url"
	"time"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/util"
)

type Track struct {
	id       id.Track
	title    string
	author   *string
	duration *time.Duration
	url      url.URL
	imageURL *url.URL
}

func New(trackID *id.Track, title string, author *string, duration *time.Duration, url url.URL, imageURL *url.URL) *Track {
	if trackID == nil {
		trackID = util.Ptr(id.GenerateTrack())
	}

	return &Track{
		id:       *trackID,
		title:    title,
		author:   author,
		duration: duration,
		url:      url,
		imageURL: imageURL,
	}
}

func (t *Track) ID() id.Track {
	return t.id
}

func (t *Track) Title() string {
	return t.title
}

func (t *Track) HasAuthor() bool {
	return t.author != nil
}

func (t *Track) Author() *string {
	return t.author
}

func (t *Track) HasDuration() bool {
	return t.duration != nil
}

func (t *Track) Duration() *time.Duration {
	return t.duration
}

func (t *Track) URL() url.URL {
	return t.url
}

func (t *Track) HasImage() bool {
	return t.imageURL != nil
}

func (t *Track) ImageURL() *url.URL {
	return t.imageURL
}
