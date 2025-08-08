package groupprovider

import (
	"context"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
)

func (p *Provider) Create(ctx context.Context, groupToCreate *group.Group) error {
	err := p.groupRepository.Create(ctx, groupToCreate)
	if err != nil {
		return fmt.Errorf("failed to create group in repository: %s", err.Error())
	}

	return nil
}
