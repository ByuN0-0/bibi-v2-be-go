package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func (c *PingCommand) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Pong! Bot의 응답 시간을 측정합니다",
	}
}

func (c *PingCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	response := fmt.Sprintf("Pong! 응답 시간: %dms", s.HeartbeatLatency().Milliseconds())

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
