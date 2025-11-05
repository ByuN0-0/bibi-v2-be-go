package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) Description() string {
	return "Pong! Bot의 응답 시간을 측정합니다"
}

func (c *PingCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	response := fmt.Sprintf("Pong! 응답 시간: %dms", s.HeartbeatLatency().Milliseconds())
	_, err := s.ChannelMessageSend(m.ChannelID, response)
	return err
}
