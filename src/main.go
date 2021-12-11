package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessageReactions

	if err != nil {
		panic(err)
	}

	session.AddHandler(messageCreate)
	session.AddHandler(messageReactionAdd)
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
