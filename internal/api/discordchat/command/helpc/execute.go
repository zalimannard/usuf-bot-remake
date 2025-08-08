package helpc

import (
	"context"
	"fmt"
	"strings"
	"usuf-bot-remake/internal/api/discordchat/command"
	"usuf-bot-remake/internal/api/discordchat/command/clearc"
	"usuf-bot-remake/internal/api/discordchat/command/loopc"
	"usuf-bot-remake/internal/api/discordchat/command/loopqc"
	"usuf-bot-remake/internal/api/discordchat/command/playc"
	"usuf-bot-remake/internal/api/discordchat/command/randomc"
	"usuf-bot-remake/internal/api/discordchat/command/skipc"
	"usuf-bot-remake/internal/domain/entity/helprow"

	"github.com/rs/zerolog/log"
)

func (c *Command) Execute(ctx context.Context, args []string) {
	if len(args) != 1 {
		log.Ctx(ctx).Error().Msg("Invalid arguments")
		return
	}

	commands := []command.Command{
		playc.New(nil),
		skipc.New(nil),
		loopc.New(nil),
		loopqc.New(nil),
		randomc.New(nil),
		clearc.New(nil),
		New(nil),
	}

	helpRows := make([]helprow.Row, 0, len(commands))

	for _, iCommand := range commands {
		mainCommand := fmt.Sprintf("%s%s ", args[0], iCommand.Names()[0])
		subcommands := ""
		for i := 1; i < len(iCommand.Names()); i++ {
			if i == 1 {
				subcommands += ""
			}
			subcommands += fmt.Sprintf("[%s%s]", args[0], iCommand.Names()[i])
		}
		parameterList := make([]string, 0)
		for i := 0; i < len(iCommand.Parameters()); i++ {
			parameterList = append(parameterList, fmt.Sprintf("(%s)", iCommand.Parameters()[i]))
		}
		title := fmt.Sprintf("%s%s %s", mainCommand, subcommands, strings.Join(parameterList, " / "))
		helpRows = append(helpRows, helprow.Row{
			Title:       title,
			Description: iCommand.Description(),
		})
	}

	err := c.helpUseCase.Help(ctx, helpRows)
	if err != nil {
		log.Ctx(ctx).Error().Err(fmt.Errorf("failed to loop track: %w", err)).Msg("Error")
	}
}
