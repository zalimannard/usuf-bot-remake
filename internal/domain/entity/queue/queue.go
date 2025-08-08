package queue

import (
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/util"
)

type Queue struct {
	queueID       id.Queue
	items         []Item
	orderType     OrderType
	currentNumber int
}

func New(queueID *id.Queue, items []Item, orderType OrderType, currentNumber int) (*Queue, error) {
	if queueID == nil {
		queueID = util.Ptr(id.GenerateQueue())
	}

	// TODO: Проверка корректности

	return &Queue{
		queueID:       *queueID,
		items:         items,
		orderType:     orderType,
		currentNumber: currentNumber,
	}, nil
}

func (q *Queue) ID() id.Queue {
	return q.queueID
}

func (q *Queue) Items() []Item {
	return q.items
}

func (q *Queue) OrderType() OrderType {
	return q.orderType
}

func (q *Queue) CurrentNumber() int {
	return q.currentNumber
}

func (q *Queue) Length() int {
	return len(q.items)
}
