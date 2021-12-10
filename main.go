package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var CDN_CHANNEL_ID = "918725182330400788"
var PICS_CHANNEL_ID = "918355152493215764"

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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	/* lol when crash */
	fmt.Println("uh oh D:")

	session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "pong")
		}

		if m.Content == "wump" {
			s.ChannelMessageSend(m.ChannelID, "<:wumpWave:918629841836859412>")
		}
	}
}
