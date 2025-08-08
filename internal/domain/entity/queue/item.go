package queue

import (
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/util"
)

type Item struct {
	itemID      id.QueueItem
	trackID     id.Track
	requesterID id.User
}

func NewItem(itemID *id.QueueItem, trackID id.Track, requesterID id.User) Item {
	if itemID == nil {
		itemID = util.Ptr(id.GenerateQueueItem())
	}

	return Item{
		itemID:      *itemID,
		trackID:     trackID,
		requesterID: requesterID,
	}
}

func (i *Item) ItemID() id.QueueItem {
	return i.itemID
}

func (i *Item) TrackID() id.Track {
	return i.trackID
}

func (i *Item) RequesterID() id.User {
	return i.requesterID
}
