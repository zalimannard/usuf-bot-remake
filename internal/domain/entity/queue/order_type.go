package queue

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownOrderType = errors.New("unknown order type")
)

func newErrUnknownOrderType(orderType string) error {
	return fmt.Errorf("%w (%s)", ErrUnknownOrderType, orderType)
}

type OrderType string

const (
	OrderTypeUnknown   OrderType = "unknown"
	OrderTypeNormal    OrderType = "normal"
	OrderTypeLoopTrack OrderType = "loop_track"
	OrderTypeLoopQueue OrderType = "loop_queue"
	OrderTypeRandom    OrderType = "random"
)

func ParseOrderType(orderType string) (OrderType, error) {
	switch orderType {
	case string(OrderTypeNormal):
		return OrderTypeNormal, nil
	case string(OrderTypeLoopTrack):
		return OrderTypeLoopTrack, nil
	case string(OrderTypeLoopQueue):
		return OrderTypeLoopQueue, nil
	case string(OrderTypeRandom):
		return OrderTypeRandom, nil
	default:
		return OrderTypeUnknown, newErrUnknownOrderType(orderType)
	}
}

func (o OrderType) String() string {
	return string(o)
}
