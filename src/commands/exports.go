package commands

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

type BotCommand struct {
	Name string
	Run  CommandHandler
}

func New(name string, handler CommandHandler) *BotCommand {
	return &BotCommand{
		Name: name,
		Run:  handler,
	}
}

var Commands = map[string]*BotCommand{
	"help": New("help", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		s.ChannelMessageSend(m.ChannelID, "This is a help message")
	}),

	"ping": New("ping", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}),

	"count": New("count", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		guildJson, err := s.RequestWithBucketID("GET", discordgo.EndpointGuild(m.GuildID)+"?with_counts=true", nil, discordgo.EndpointGuild(m.GuildID))
		guildDiscord := &discordgo.Guild{}

		if err == nil {
			_ = json.Unmarshal(guildJson, guildDiscord)
		}

		// Figure out later

		// count := 0

		// for _, member := range guildDiscord.Members {
		// 	if utils.Contains(member.Roles, utils.OWNER_ROLE_ID) {
		// 		count += 1
		// 		s.ChannelMessageSend(m.ChannelID, fmt.Sprint(member))
		// 	}
		// }

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Type:        "rich",
			Color:       7506394,
			Description: "**" + fmt.Sprint(guildDiscord.ApproximateMemberCount) + "** Wumpus Plushie owners currently reside in this server",
		})
	}),
}
