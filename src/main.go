package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"wumpus/src/interactions"
	"wumpus/src/utils"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildIntegrations

	if err != nil {
		panic(err)
	}

	if utils.APP_ENV == "production" {
		session.AddHandler(guildMemberAdd)
		session.AddHandler(guildMemberUpdate)
		session.AddHandler(messageReactionAdd)

		// Slash commands
		session.AddHandler(interactions.InteractionReceived)
	}

	session.AddHandler(messageCreate)

	fmt.Printf("Running in %s mode\n", utils.APP_ENV)

	session.Open()

	if utils.APP_ENV == "production" && session.State.User.ID == utils.PROD_BOT_ID {
		interactions.CreateCommands(session)
	}

	if session.State.User.ID == utils.PROD_BOT_ID {
		session.UpdateGameStatus(0, "big wumpus")
	} else {
		session.UpdateGameStatus(0, "small wumpus")
	}

	fmt.Println("🚀 Wumpus has launched :D")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	fmt.Println("uh oh D:")

	session.Close()
}
