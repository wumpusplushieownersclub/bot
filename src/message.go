package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"wumpus/src/cdn"
	"wumpus/src/commands"
	"wumpus/src/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/dgrijalva/jwt-go"
)

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.MessageReaction.UserID == s.State.User.ID {
		return
	}

	member, _ := s.GuildMember(r.GuildID, r.MessageReaction.UserID)

	if r.ChannelID == utils.VERIFICATION_CHANNEL_ID {
		if utils.Contains(member.Roles, utils.TEAM_ROLE_ID) && utils.Contains(utils.VALID_REACTIONS, r.Emoji.Name) {
			reactedMessage, _ := s.ChannelMessage(r.ChannelID, r.MessageID)
			channelMessages, _ := s.ChannelMessages(r.ChannelID, 100, r.MessageID, "", "")
			var originalMessage *discordgo.Message

			for _, v := range channelMessages {
				if v.Author.ID == s.State.User.ID {
					for _, m := range v.Mentions {
						if m.ID == reactedMessage.Author.ID {
							originalMessage = v
							s.ChannelMessageDelete(r.ChannelID, originalMessage.ID)
						}
					}
				}
			}

			image := reactedMessage.Attachments[0]
			httpClient := &http.Client{
				Timeout: 10 * time.Second,
			}
			request, _ := httpClient.Get(image.URL)
			nameSplit := strings.Split(image.Filename, ".")

			filename := fmt.Sprintf("%s-%s.%s", reactedMessage.Author.ID, reactedMessage.ID, nameSplit[len(nameSplit)-1])
			mime := utils.EXT_TO_MIME[nameSplit[len(nameSplit)-1]]

			url, uploadErr := cdn.UploadToCdn(request.Body, "verifications", filename, mime)
			if uploadErr != nil {
				fmt.Println("There as an error uploading to CDN", uploadErr)
			}

			s.ChannelMessageDelete(r.ChannelID, reactedMessage.ID)

			var status string
			var color int

			if r.Emoji.Name == "👍" {
				status = "verified"
				color = 0x00FF00
				s.GuildMemberRoleAdd(r.GuildID, reactedMessage.Author.ID, utils.OWNER_ROLE_ID)
			} else if r.Emoji.Name == "👎" {
				status = "denied"
				color = 0xFF0000
				s.GuildMemberDelete(r.GuildID, reactedMessage.Author.ID)
			}

			message := fmt.Sprintf("%s#%s was %s by %s", reactedMessage.Author.Username, reactedMessage.Author.Discriminator, status, member.User.Username)

			s.ChannelMessageSendEmbed(utils.LOGS_CHANNEL_ID, &discordgo.MessageEmbed{
				Type:      "rich",
				Color:     color,
				Title:     message,
				Image:     &discordgo.MessageEmbedImage{URL: url},
				Footer:    &discordgo.MessageEmbedFooter{Text: "Wumpus Verification"},
				Timestamp: time.Now().Format(time.RFC3339),
			})
		} else {
			s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.MessageReaction.UserID)
		}
	} else if r.ChannelID == utils.PICS_CHANNEL_ID {
		if r.Emoji.Name == "🔥" {
			reactedMessage, _ := s.ChannelMessage(r.ChannelID, r.MessageID)
			timestamp, _ := reactedMessage.Timestamp.Parse()

			if int(timestamp.Unix()) < int(1639342558) {
				return
			}

			// Get the users id from the image file path (we set this so it's trustable)
			filename := strings.Split(reactedMessage.Embeds[0].Image.URL, "/")
			splits := strings.Split(filename[len(filename)-1], ".")
			originalAuthor, _ := s.GuildMember(r.GuildID, splits[0])

			if originalAuthor.User.ID == member.User.ID {
				return
			}

			go func() {
				if utils.POINTS_WORKER_SECRET != "provide_in_env" {
					httpClient := &http.Client{
						Timeout: 10 * time.Second,
					}
					jsonPayload := []byte(fmt.Sprintf(`{
						"message_id": %s
					}`, reactedMessage.ID))
					claims := jwt.MapClaims{}
					claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
					withClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					token, _ := withClaims.SignedString([]byte(utils.POINTS_WORKER_SECRET))
					request, _ := http.NewRequest("POST", fmt.Sprintf("%s/track/firepic/%s", utils.POINTS_WORKER_HOST, originalAuthor.User.ID), bytes.NewBuffer(jsonPayload))
					request.Header.Set("content-type", "application/json")
					request.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))
					httpClient.Do(request)
				}
			}()
		}
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID == utils.PICS_CHANNEL_ID && utils.APP_ENV == "production" {
		if len(m.Attachments) > 0 && m.Attachments[0].Height != 0 && m.Attachments[0].Width != 0 && !utils.IsVideo(m.Attachments[0]) {
			image := m.Attachments[0]
			httpClient := &http.Client{
				Timeout: 10 * time.Second,
			}
			request, _ := httpClient.Get(image.URL)
			nameSplit := strings.Split(image.Filename, ".")

			filename := fmt.Sprintf("%s-%s.%s", m.Author.ID, m.ID, nameSplit[len(nameSplit)-1])
			mime := utils.EXT_TO_MIME[nameSplit[len(nameSplit)-1]]

			url, uploadErr := cdn.UploadToCdn(request.Body, "pics", filename, mime)
			if uploadErr != nil {
				fmt.Println("There as an error uploading to cdn", uploadErr)
				errMessage, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Type:        "rich",
					Color:       0x7289DA,
					Description: "Error uploading file to wumpus cdn, report this",
				})
				time.AfterFunc(5*time.Second, func() { s.ChannelMessageDelete(m.ChannelID, errMessage.ID) })
				return
			}

			message, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       2617723,
				Image:       &discordgo.MessageEmbedImage{URL: url},
				Author:      &discordgo.MessageEmbedAuthor{IconURL: m.Author.AvatarURL(""), Name: fmt.Sprintf("%s#%s", m.Author.Username, m.Author.Discriminator)},
				Description: m.Content,
			})

			s.MessageReactionAdd(m.ChannelID, message.ID, "🔥")
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		} else {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			message, _ := s.ChannelMessageSend(m.ChannelID, "Please include an image!")
			time.AfterFunc(5*time.Second, func() { s.ChannelMessageDelete(m.ChannelID, message.ID) })
		}
	} else if m.ChannelID == utils.VERIFICATION_CHANNEL_ID && utils.APP_ENV == "production" {
		if len(m.Attachments) > 0 && m.Attachments[0].Height != 0 && m.Attachments[0].Width != 0 {
			s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
			s.MessageReactionAdd(m.ChannelID, m.ID, "👎")
		} else {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			message, _ := s.ChannelMessageSend(m.ChannelID, "Must include an image to get verified!")
			time.AfterFunc(5*time.Second, func() { s.ChannelMessageDelete(m.ChannelID, message.ID) })
		}

		return
	}

	if (strings.ToLower(m.Content) == "wump" && utils.APP_ENV == "production" && s.State.User.ID == utils.PROD_BOT_ID) || (m.ContentWithMentionsReplaced() == fmt.Sprintf("@%s", s.State.User.Username) && m.Mentions[0].ID == s.State.User.ID && utils.APP_ENV == "production") {
		if m.ChannelID == utils.PICS_CHANNEL_ID || m.ChannelID == utils.VERIFICATION_CHANNEL_ID {
			return
		}

		s.ChannelMessageSend(m.ChannelID, "<:wumpWave:918629841836859412>")
	} else {
		if m.ChannelID == utils.PICS_CHANNEL_ID || m.ChannelID == utils.VERIFICATION_CHANNEL_ID || (utils.APP_ENV != "production" && utils.PREFIX == utils.DEFAULT_PREFIX) {
			return
		}

		contentLower := strings.ToLower(m.Content)

		if !strings.HasPrefix(contentLower, utils.PREFIX) {
			return
		}

		contentSplit := strings.Split(m.Content[len(utils.PREFIX):], " ")
		commandName := strings.ToLower(contentSplit[0])
		command := commands.Commands[commandName]

		if command == nil {
			s.ChannelMessageSend(m.ChannelID, "Unknown command")
			return
		}

		args := contentSplit[1:]

		go command.Run(s, m, args)
	}

	if utils.POINTS_WORKER_SECRET != "provide_in_env" && (utils.APP_ENV == "production" || utils.POINTS_WORKER_HOST != utils.DEFAULT_POINTS_WORKER_HOST) {
		go func() {
			httpClient := &http.Client{
				Timeout: 10 * time.Second,
			}
			claims := jwt.MapClaims{}
			claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
			withClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token, _ := withClaims.SignedString([]byte(utils.POINTS_WORKER_SECRET))
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/track/message/%s", utils.POINTS_WORKER_HOST, m.Author.ID), nil)
			request.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))
			httpClient.Do(request)
		}()
	}
}
