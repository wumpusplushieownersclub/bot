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

// Create environment variables locally for base development
var CDN_CHANNEL_ID = getenv("CDN_CHANNEL", "918725182330400788")
var PICS_CHANNEL_ID = getenv("PICS_CHANNEL", "918355152493215764")
var TEAM_ROLE_ID = getenv("TEAM_ROLE", "918354701337116703")

func main() {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	session.Identify.Intents = discordgo.IntentsGuildMessages

	// If there is an error, print it and return.
	if err != nil {
		panic(err)
	}

	session.AddHandler(messageCreate)

	// Connect to the Discord API with the token provided -- Session.Open();
	session.Open()

	/* When bot starts up -- Log to the console */
	fmt.Println("🚀 Bot has launched")

	session.UpdateGameStatus(0, "big wumpus")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	/* When crash -> Print Error */
	fmt.Println("🛑 Uh oh! It appears that an error has occured.")

	session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// If the user is the bot, return the function.
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

				// If image is provided, embed it
				Type:        "rich",
				Color:       2617723,
				Image:       &discordgo.MessageEmbedImage{URL: cdnMessage.Attachments[0].URL},
				Author:      &discordgo.MessageEmbedAuthor{IconURL: m.Author.AvatarURL(""), Name: m.Author.Username},
				Description: m.Content,
			})
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		} else {
			// If it doesn't have an attachment of an image return an error
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
