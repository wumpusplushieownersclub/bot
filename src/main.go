package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var CDN_CHANNEL_ID = getenv("CDN_CHANNEL", "918725182330400788")
var PICS_CHANNEL_ID = getenv("PICS_CHANNEL", "918355152493215764")
var LOGS_CHANNEL_ID = getenv("LOGS_CHANNEL", "918952346975862824")
var VERIFICATION_CHANNEL_ID = getenv("VERIFICATION_CHANNEL", "918932836428419163")

var TEAM_ROLE_ID = getenv("TEAM_ROLE", "918354701337116703")
var OWNER_ROLE_ID = getenv("OWNER_ROLE", "918355466894065685")

var VALID_REACTIONS = []string{"👍", "👎"}

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
