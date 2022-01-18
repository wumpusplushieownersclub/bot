package utils

import (
	"strings"

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

func IsVideo(attachment *discordgo.MessageAttachment) bool {
	nameSplit := strings.Split(attachment.Filename, ".")
	ext := nameSplit[len(nameSplit)-1]

	for _, t := range VIDEO_FORMATS {
		if strings.ToLower(ext) == t {
			return true
		}
	}

	return false
}
