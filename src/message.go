package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageReaction.UserID == s.State.User.ID || r.ChannelID != VERIFICATION_CHANNEL_ID {
		return
	}

	member, _ := s.GuildMember(r.GuildID, r.MessageReaction.UserID)

	if contains(member.Roles, TEAM_ROLE_ID) && contains(VALID_REACTIONS, r.Emoji.Name) {
		reactedMessage, _ := s.ChannelMessage(r.ChannelID, r.MessageID)
		channelMessages, _ := s.ChannelMessages(r.ChannelID, 100, r.MessageID, "", "")
		var originalMessage *discordgo.Message

		for _, v := range channelMessages {
			if v.Author.ID == s.State.User.ID {
				for _, m := range v.Mentions {
					if m.ID == reactedMessage.Author.ID {
						originalMessage = v
					}
				}
			}
		}

		s.ChannelMessageDelete(r.ChannelID, reactedMessage.ID)
		s.ChannelMessageDelete(r.ChannelID, originalMessage.ID)

		var status string
		var color int

		if r.Emoji.Name == "👍" {
			status = "verified"
			color = 0x00FF00
			s.GuildMemberRoleAdd(r.GuildID, reactedMessage.Author.ID, OWNER_ROLE_ID)
		} else if r.Emoji.Name == "👎" {
			status = "denied"
			color = 0xFF0000
			s.GuildMemberDelete(r.GuildID, reactedMessage.Author.ID)
		}

		message := fmt.Sprintf("%s#%s was %s by %s", reactedMessage.Author.Username, reactedMessage.Author.Discriminator, status, member.User.Username)

		s.ChannelMessageSendEmbed(LOGS_CHANNEL_ID, &discordgo.MessageEmbed{
			Type:      "rich",
			Color:     color,
			Title:     message,
			Footer:    &discordgo.MessageEmbedFooter{Text: "Wumpus Verification"},
			Timestamp: time.Now().Format(time.RFC3339),
		})
	} else {
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.MessageReaction.UserID)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID == PICS_CHANNEL_ID {
		if len(m.Attachments) > 0 && m.Attachments[0].Height != 0 && m.Attachments[0].Width != 0 {
			image := m.Attachments[0]
			httpClient := &http.Client{
				Timeout: 10 * time.Second,
			}
			request, _ := httpClient.Get(image.URL)

			cdnMessage, _ := s.ChannelFileSend(CDN_CHANNEL_ID, image.Filename, request.Body)

			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       2617723,
				Image:       &discordgo.MessageEmbedImage{URL: cdnMessage.Attachments[0].URL},
				Author:      &discordgo.MessageEmbedAuthor{IconURL: m.Author.AvatarURL(""), Name: m.Author.Username},
				Description: m.Content,
			})

			s.ChannelMessageDelete(m.ChannelID, m.ID)
		} else {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			message, _ := s.ChannelMessageSend(m.ChannelID, "Please include an image!")
			time.AfterFunc(5*time.Second, func() { s.ChannelMessageDelete(m.ChannelID, message.ID) })
		}
	} else if m.ChannelID == VERIFICATION_CHANNEL_ID {
		if len(m.Attachments) > 0 && m.Attachments[0].Height != 0 && m.Attachments[0].Width != 0 {
			s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
			s.MessageReactionAdd(m.ChannelID, m.ID, "👎")
		} else {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			message, _ := s.ChannelMessageSend(m.ChannelID, "Must include an image to get verified!")
			time.AfterFunc(5*time.Second, func() { s.ChannelMessageDelete(m.ChannelID, message.ID) })
		}
	} else {
		if strings.ToLower(m.Content) == "ping" {
			s.ChannelMessageSend(m.ChannelID, "pong")
		}

		if strings.ToLower(m.Content) == "wump" {
			s.ChannelMessageSend(m.ChannelID, "<:wumpWave:918629841836859412>")
		}

		if strings.ToLower(m.Content) == "nap" && contains(m.Member.Roles, TEAM_ROLE_ID) {
			s.ChannelMessageSend(m.ChannelID, "<:wumpSad:918629842050748437> going down for nap time")
			s.Close()
			os.Exit(9)
		}

		if strings.ToLower(m.Content) == "wump count" {
			guildJson, err := s.RequestWithBucketID("GET", discordgo.EndpointGuild(m.GuildID)+"?with_counts=true", nil, discordgo.EndpointGuild(m.GuildID))
            guildDiscord := &discordgo.Guild{}

            if err == nil {
				_ = json.Unmarshal(guildJson, guildDiscord)
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       7506394,
				Description: "**" + fmt.Sprint(guildDiscord.ApproximateMemberCount) + "** Wumpus Plushie owners currently reside in this server",
			})
		}
	}
}
