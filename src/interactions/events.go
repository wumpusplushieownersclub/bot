package interactions

import (
	"github.com/bwmarrin/discordgo"
)

func InteractionReceived(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, slashCommand := range Commands {
		if slashCommand.Command.Name == i.ApplicationCommandData().Name {
			slashCommand.Run(s, i)
		}
	}
}
