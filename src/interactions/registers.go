package interactions

import (
	"fmt"
	"strconv"
	"wumpus/src/points"
	"wumpus/src/utils"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type SlashCommand struct {
	Command discordgo.ApplicationCommand
	Run     CommandHandler
}

func DeleteCommands(s *discordgo.Session) {
	commands, _ := s.ApplicationCommands(s.State.User.ID, utils.GUILD_ID)
	for _, ac := range commands {
		s.ApplicationCommandDelete(s.State.User.ID, utils.GUILD_ID, ac.ID)
	}
}

func CreateCommands(s *discordgo.Session) {
	for _, slashCommand := range Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, utils.GUILD_ID, &slashCommand.Command)

		if err != nil {
			fmt.Println("Failed to create command", slashCommand.Command.Name, err)
			return
		}

		fmt.Println("Create command", cmd.Name)
	}
}

func New(command discordgo.ApplicationCommand, handler CommandHandler) *SlashCommand {
	return &SlashCommand{
		Command: command,
		Run:     handler,
	}
}

var Commands = map[string]*SlashCommand{
	"balance": New(discordgo.ApplicationCommand{
		Name:        "balance",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Get WumpCoin balance",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		embeds := make([]*discordgo.MessageEmbed, 0)

		pointsAccount, pointsErr := points.GetPointsAccount(i.Member.User.ID)
		if pointsErr != nil {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       0x7289DA,
				Description: "Failed to get balance from the server",
			})
		} else {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       0x7289DA,
				Description: fmt.Sprintf("You have **%s** WumpCoins", strconv.Itoa(int(pointsAccount.Points))),
			})
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embeds,
				Flags:  1 << 6,
			},
		})
	}),

	"leaderboard": New(discordgo.ApplicationCommand{
		Name:        "leaderboard",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Get Top 10 WumpCoin leaderboard",
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		embeds := make([]*discordgo.MessageEmbed, 0)

		board, err := points.GetLeaderboard()

		if err != nil {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Type:        "rich",
				Color:       0x7289DA,
				Description: "Failed to get leaderboard from the server",
			})

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: embeds,
					Flags:  1 << 6,
				},
			})

			return
		}

		leadboard := ""

		for index, account := range board {
			member, _ := s.GuildMember(utils.GUILD_ID, account.User)
			leadboard += fmt.Sprintf("**%x**. %s#%s : **%s**\n", index+1, member.User.Username, member.User.Discriminator, strconv.Itoa(int(account.Points)))
		}

		embeds = append(embeds, &discordgo.MessageEmbed{
			Type:        "rich",
			Color:       0x7289DA,
			Title:       "WumpCoin Leaderboard",
			Description: leadboard,
		})

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embeds,
				Flags:  1 << 6,
			},
		})
	}),
}
