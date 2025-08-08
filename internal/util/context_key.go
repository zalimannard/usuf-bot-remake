package util

type ContextKey string

const (
	ContextKeyLogger         ContextKey = "logger"
	ContextKeyRequesterGroup            = ContextKey("requester_group")
	ContextKeyRequesterUser             = ContextKey("requester_user")
)
