package utils

import (
	"github.com/bwmarrin/discordgo"
)

func CountRoleMembers(s *discordgo.Session, guildID string, roleID string) int {
	s.RequestGuildMembers(guildID, "", 0, false)
	guildMembers, _ := s.GuildMembers(guildID, "", 1000) // This'll work for <1000 members

	count := 0
	for _, member := range guildMembers {
		if Contains(member.Roles, roleID) {
			count += 1
		}
	}

	return count
}