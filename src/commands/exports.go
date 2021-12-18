package commands

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"wumpus/src/points"
	"wumpus/src/utils"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

type BotCommand struct {
	Name        string
	Description string
	Run         CommandHandler
}

func New(name string, desc string, handler CommandHandler) *BotCommand {
	return &BotCommand{
		Name:        name,
		Description: desc,
		Run:         handler,
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
		count := utils.CountRoleMembers(s, m.GuildID, utils.OWNER_ROLE_ID)

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Type:        "rich",
			Color:       0x7289DA,
			Description: "**" + fmt.Sprint(count) + "** Wumpus Plushie owners currently reside in this server",
		})
	}),

	"nap": New("nap", "bye bye wumpus", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		if utils.Contains(m.Member.Roles, utils.TEAM_ROLE_ID) {
			if utils.APP_ENV == "production" {
				s.ChannelMessageSend(m.ChannelID, "<:wumpAngry:918630551995432960> How dare you try to make me nap, it's work time baby")
				return
			}

			s.ChannelMessageSend(m.ChannelID, "<:wumpSad:918629842050748437> going down for nap time")
			s.Close()
			os.Exit(9)
		} else {
			s.ChannelMessageSend(m.ChannelID, "<:wumpSad:918629842050748437> Only the team can make me nap")
		}
	}),

	"balance": New("balance", "Get your WumpCoin balance", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		pointsAccount, pointsErr := points.GetPointsAccount(m.Author.ID)

		if pointsErr != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       0x7289DA,
				Description: "Failed to get balance from the server",
			})
			return
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Type:        "rich",
			Color:       0x7289DA,
			Description: fmt.Sprintf("You have **%s** WumpCoins", strconv.Itoa(int(pointsAccount.Points))),
			Author:      &discordgo.MessageEmbedAuthor{IconURL: m.Author.AvatarURL(""), Name: fmt.Sprintf("%s#%s", m.Author.Username, m.Author.Discriminator)},
		})
	}),

	"leaderboard": New("leaderboard", "See top WumpCoin holders", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		board, err := points.GetLeaderboard()

		if err != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       0x7289DA,
				Description: "Failed to get leaderboard from the server",
			})
			return
		}

		leadboard := ""

		for index, account := range board {
			member, _ := s.GuildMember(m.GuildID, account.User)
			leadboard += fmt.Sprintf("**%x**. %s#%s : **%s**\n", index+1, member.User.Username, member.User.Discriminator, strconv.Itoa(int(account.Points)))
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Type:        "rich",
			Color:       0x7289DA,
			Title:       "WumpCoin Leaderboard",
			Description: leadboard,
		})
	}),

	"purge": New("purge", "Purge messages in a channel from 2-100", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		if utils.Contains(m.Member.Roles, utils.TEAM_ROLE_ID) {
			count, _ := strconv.Atoi(args[0])

			if count < 2 || count > 100 {
				msg, _ := s.ChannelMessageSend(m.ChannelID, "<:wumpSad:918629842050748437> Must include a range between 2-100 to purge messages")
				time.AfterFunc(3*time.Second, func() {
					s.ChannelMessageDelete(m.ChannelID, msg.ID)
				})
				return
			}

			channelMessages, _ := s.ChannelMessages(m.ChannelID, count, m.ID, "", "")
			messages := make([]string, 0)
			for _, message := range channelMessages {
				messages = append(messages, message.ID)
			}
			s.ChannelMessagesBulkDelete(m.ChannelID, messages)
			msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Purged **%s** messages", strconv.Itoa(len(messages))))
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			time.AfterFunc(5*time.Second, func() {
				s.ChannelMessageDelete(m.ChannelID, msg.ID)
			})
		} else {
			s.ChannelMessageSend(m.ChannelID, "<:wumpSad:918629842050748437> Only the team can purge messages")
		}
	}),

	"info": New("info", "Display user information", func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		lookup := m.Author
		member := m.Member

		if len(m.Mentions) > 0 {
			lookup = m.Mentions[0]
			member, _ = s.GuildMember(m.GuildID, lookup.ID)
		} else if len(args) > 0 {
			mem, err := s.GuildMember(m.GuildID, args[0])
			lookup = mem.User

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "That user doesn't appear to be in this server... <:wumpSad:918629842050748437>")
			} else {
				member = mem
			}
		}

		joinedTime, _ := member.JoinedAt.Parse()

		pointsAccount, pointsErr := points.GetPointsAccount(lookup.ID)

		if pointsErr != nil {
			pointsAccount = &points.PointsAccount{User: lookup.ID, Points: 0}
		}

		fields := make([]*discordgo.MessageEmbedField, 0)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "User ID",
			Value:  fmt.Sprintf("`%s`", lookup.ID),
			Inline: false,
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Joined At",
			Value:  fmt.Sprintf("`%s`", joinedTime.Format("01-02-2006 15:04:05")),
			Inline: false,
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "WumpCoins",
			Value:  fmt.Sprintf("**%s**", strconv.Itoa(int(pointsAccount.Points))),
			Inline: false,
		})

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Type:      "rich",
			Color:     0x7289DA,
			Title:     fmt.Sprintf("%s#%s", lookup.Username, lookup.Discriminator),
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: lookup.AvatarURL("512")},
			Fields:    fields,
		})
	}),
}
