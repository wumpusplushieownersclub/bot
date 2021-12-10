package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

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

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}

	if m.Content == "wump" {
		s.ChannelMessageSend(m.ChannelID, "<:wumpWave:918629841836859412>")
	}

	/* change id to a variable so we can change easier or something idk */
	if m.ChannelID == "918355152493215764" {
		if len(m.Attachments) > 0 {
			image := m.Attachments[0]
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type: "rich", 
				Image: &discordgo.MessageEmbedImage{URL: image.URL}, 
				Author: &discordgo.MessageEmbedAuthor{IconURL: m.Author.AvatarURL(""), Name: m.Author.Username}, 
				Description: m.Content,
			})
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Please include an image!")
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	}
}
