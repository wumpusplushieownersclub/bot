package main

import (
	"wumpus/src/utils"

	"github.com/bwmarrin/discordgo"
)

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend(utils.VERIFICATION_CHANNEL_ID, "Hey "+m.Mention()+", in order to get access to the community, please send an image of your Wumpus plushie next to a note containg your Discord username and tag in this channel.")
}
