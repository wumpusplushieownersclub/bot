package commands

import (
	"fmt"
	"wumpus/src/utils"

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
		s.RequestGuildMembers(m.GuildID, "", 0, false)
		guildMembers, _ := s.GuildMembers(m.GuildID, "", 1000) // This'll work for <1000 members

		count := 0
		for _, member := range guildMembers {
			if utils.Contains(member.Roles, utils.OWNER_ROLE_ID) {
				count += 1
			}
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Type:        "rich",
			Color:       7506394,
			Description: "**" + fmt.Sprint(count) + "** Wumpus Plushie owners currently reside in this server",
		})
	}),
}
