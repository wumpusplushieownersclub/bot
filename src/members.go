package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	fmt.Println("Member joined", m.User.Username)
}
