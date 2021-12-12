package commands

import (
	"fmt"
	"os"
	"wumpus/src/utils"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

type BotCommand struct {
	Name string
	Description string
	Run  CommandHandler
}

func New(name string, desc string, handler CommandHandler) *BotCommand {
	return &BotCommand{
		Name: name,
		Description: desc,
		Run:  handler,
	}
}

var Commands = map[string]*BotCommand{
	"help": New("help", "Display help information", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

		// I cannot access map information from inside the map itself, need to figure out how to make it dynamic

		// commandInfo := ""

		// for name, command := range Commands {
		// 	commandInfo += fmt.Sprintf("`%s`: %s\n", name, command.Description)
		// }

		// s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		// 	Type:        "rich",
		// 	Color:       0x4E5D94,
		// 	Description: fmt.Sprintf("**Wumpus Commands**\n%s", commandInfo),
		// })

		// For now, use this

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Type:        "rich",
			Color:       0x4E5D94,
			Description: "**Wumpus Commands (prefixed with `wump`)**\n\n`help`: Display help information\n`ping`: pong :D\n`count`: Return how many Wumpus Plushie owners are in the guild",
		})

	}),

	"ping": New("ping", "pong :D", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}),

	"count": New("count", "Return how many Wumpus Plushie owners are in the guild.", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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
			Color:       0x7289DA,
			Description: "**" + fmt.Sprint(count) + "** Wumpus Plushie owners currently reside in this server",
		})
	}),

	"nap": New("nap", "bye bye wumpus", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		if utils.Contains(m.Member.Roles, utils.TEAM_ROLE_ID) {
			s.ChannelMessageSend(m.ChannelID, "<:wumpSad:918629842050748437> going down for nap time")
			s.Close()
			os.Exit(9)
		}
	}),
}
