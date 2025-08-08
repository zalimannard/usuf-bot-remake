package middleware

import (
	"context"
	"errors"
	"fmt"
	"usuf-bot-remake/internal/domain/entity/group"
	"usuf-bot-remake/internal/domain/entity/id"
	"usuf-bot-remake/internal/domain/entity/user"
	"usuf-bot-remake/internal/provider/groupprovider"
	"usuf-bot-remake/internal/provider/userprovider"
	"usuf-bot-remake/internal/util"

	"github.com/bwmarrin/discordgo"
)

func (m *Middleware) RequesterInfo(ctx context.Context, message *discordgo.MessageCreate) (context.Context, error) {
	requesterGroup, err := m.groupUseCase.GetByExternalID(ctx, id.ParseGroupExternal(message.GuildID))
	if err != nil {
		if errors.Is(err, groupprovider.ErrGroupByExternalIDNotFound) {
			newGroup := group.New(nil, id.ParseGroupExternal(message.GuildID))
			err = m.groupUseCase.Create(ctx, newGroup)
			if err != nil {
				return nil, fmt.Errorf("failed to create group: %w", err)
			}
			requesterGroup = newGroup
		} else {
			return nil, fmt.Errorf("failed to get group: %w", err)
		}
	}

	requesterUser, err := m.userUseCase.GetByExternalID(ctx, id.ParseUserExternal(message.Author.ID))
	if err != nil {
		if errors.Is(err, userprovider.ErrUserByExternalIDNotFound) {
			newUser := user.New(nil, id.ParseUserExternal(message.Author.ID), message.Author.GlobalName)
			err = m.userUseCase.Create(ctx, newUser)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
			requesterUser = newUser
		} else {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	if requesterUser.Name() != message.Author.GlobalName {
		userToUpdate := user.New(util.Ptr(requesterUser.ID()), id.ParseUserExternal(message.Author.ID), message.Author.GlobalName)
		err = m.userUseCase.Update(ctx, userToUpdate)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		requesterUser = userToUpdate
	}

	ctxWithGroup := context.WithValue(ctx, util.ContextKeyRequesterGroup, requesterGroup)
	ctxWithUser := context.WithValue(ctxWithGroup, util.ContextKeyRequesterUser, requesterUser)

	return ctxWithUser, nil
}
