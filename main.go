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

var CDN_CHANNEL_ID = getenv("CDN_CHANNEL", "918725182330400788")
var PICS_CHANNEL_ID = getenv("PICS_CHANNEL", "918355152493215764")
var TEAM_ROLE_ID = getenv("TEAM_ROLE", "918354701337116703")

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	session.Identify.Intents = discordgo.IntentsGuildMessages

	if err != nil {
		panic(err)
	}

	session.AddHandler(messageCreate)

	session.Open()

	/* Log when the bot is online */
	fmt.Println("Bot is online :D")

	session.UpdateGameStatus(0, "big wumpus")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	/* lol when crash */
	fmt.Println("uh oh D:")

	session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore ourselves
	if m.Author.ID == s.State.User.ID {
		return
	}

	// This if statement prevents people from running commands in the pics channel
	if m.ChannelID == PICS_CHANNEL_ID {
		if len(m.Attachments) > 0 {
			image := m.Attachments[0]

			// Fetch attachment
			request, _ := http.Get(image.URL)

			// Send attachment in the CDN channel
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
