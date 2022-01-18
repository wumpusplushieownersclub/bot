package main

import (
	"wumpus/src/utils"

	"github.com/bwmarrin/discordgo"
)

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend(utils.VERIFICATION_CHANNEL_ID, "Hey "+m.Mention()+", in order to get access to the community, please send an image of your Wumpus plushie next to a note containing your Discord username and tag in this channel.")
}

func guildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	channelMessages, _ := s.ChannelMessages(utils.VERIFICATION_CHANNEL_ID, 100, "", "", "")
	var originalMessage *discordgo.Message

	for _, v := range channelMessages {
		if v.Author.ID == s.State.User.ID {
			for _, mention := range v.Mentions {
				if mention.ID == m.User.ID {
					originalMessage = v
					s.ChannelMessageDelete(utils.VERIFICATION_CHANNEL_ID, originalMessage.ID)
				}
			}
		}
	}
}

func guildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	// Handle update of member roles later for something
}
