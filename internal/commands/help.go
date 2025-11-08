package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type HelpCommand struct{}

func (c *HelpCommand) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "사용 가능한 명령어 목록을 표시합니다",
	}
}

func (c *HelpCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	embed := &discordgo.MessageEmbed{
		Title:       "bibi-bot 명령어 도움말",
		Description: "사용 가능한 모든 명령어를 확인하세요",
		Color:       0x00FF00,
		Fields:      []*discordgo.MessageEmbedField{},
	}

	for _, cmd := range GetAllCommands() {
		appCmd := cmd.ApplicationCommand()
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("/%s", appCmd.Name),
			Value: appCmd.Description,
		})
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
