package command

import "context"

type Command interface {
	Names() []string
	Parameters() []string
	Description() string
	Execute(ctx context.Context, args []string)
}
