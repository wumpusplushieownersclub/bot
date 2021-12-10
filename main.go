package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

// Channel variables
var CDN_CHANNEL_ID = getenv("CDN_CHANNEL", "918725182330400788")
var PICS_CHANNEL_ID = getenv("PICS_CHANNEL", "918355152493215764")
var VERIFICATION_CHANNEL_ID = getenv("VERIFICATION_CHANNEL", "918932836428419163")

// Role variables
var TEAM_ROLE_ID = getenv("TEAM_ROLE", "918354701337116703")

func main() {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	if err != nil {
		panic(err)
	}

	session.AddHandler(messageCreate)
	session.AddHandler(guildMemberAdd)

	session.Open()

	fmt.Println("🚀 Wumpus has launched :D")

	session.UpdateGameStatus(0, "big wumpus")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	fmt.Println("uh oh D:")

	session.Close()
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	fmt.Println("Member joined", m.User.Username)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID == PICS_CHANNEL_ID {
		if len(m.Attachments) > 0 && m.Attachments[0].Height != 0 && m.Attachments[0].Width != 0 {
			image := m.Attachments[0]
			request, _ := http.Get(image.URL)

			cdnMessage, _ := s.ChannelFileSend(CDN_CHANNEL_ID, image.Filename, request.Body)

			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       2617723,
				Image:       &discordgo.MessageEmbedImage{URL: cdnMessage.Attachments[0].URL},
				Author:      &discordgo.MessageEmbedAuthor{IconURL: m.Author.AvatarURL(""), Name: m.Author.Username},
				Description: m.Content,
			})

			s.ChannelMessageDelete(m.ChannelID, m.ID)
		} else {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			message, _ := s.ChannelMessageSend(m.ChannelID, "Please include an image!")
			time.AfterFunc(5*time.Second, func() { s.ChannelMessageDelete(m.ChannelID, message.ID) })
		}
	} else if m.ChannelID == VERIFICATION_CHANNEL_ID {
		if len(m.Attachments) > 0 && m.Attachments[0].Height != 0 && m.Attachments[0].Width != 0 {
			s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
			s.MessageReactionAdd(m.ChannelID, m.ID, "👎")
		} else {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			message, _ := s.ChannelMessageSend(m.ChannelID, "Must include an image to get verified!")
			time.AfterFunc(5*time.Second, func() { s.ChannelMessageDelete(m.ChannelID, message.ID) })
		}
	} else {
		if strings.ToLower(m.Content) == "ping" {
			s.ChannelMessageSend(m.ChannelID, "pong")
		}

		if strings.ToLower(m.Content) == "wump" {
			s.ChannelMessageSend(m.ChannelID, "<:wumpWave:918629841836859412>")
		}

		if strings.ToLower(m.Content) == "nap" && contains(m.Member.Roles, TEAM_ROLE_ID) {
			s.ChannelMessageSend(m.ChannelID, "<:wumpSad:918629842050748437> going down for nap time")
			s.Close()
			os.Exit(9)
		}

	}
}
