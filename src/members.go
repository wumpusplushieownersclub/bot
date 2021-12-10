package main

import (
	"github.com/bwmarrin/discordgo"
)

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend(VERIFICATION_CHANNEL_ID, "Hey "+m.Mention()+" in order to get access to the community, you'll need to send an image with your wumpus plushie next to a note containg your discord tag in this channel.")
}
