package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Command 슬래쉬 명령어 인터페이스
type Command interface {
	// ApplicationCommand Discord에 등록할 ApplicationCommand를 반환합니다
	ApplicationCommand() *discordgo.ApplicationCommand

	// Execute 명령어를 실행합니다 (Interaction 기반)
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

// CommandRegistry 등록된 명령어들
var commandRegistry = make(map[string]Command)

// Register 명령어를 등록합니다
func Register(cmd Command) {
	commandRegistry[cmd.ApplicationCommand().Name] = cmd
}

// GetCommand 명령어를 조회합니다
func GetCommand(name string) Command {
	return commandRegistry[name]
}

// GetAllCommands 모든 명령어를 반환합니다
func GetAllCommands() []Command {
	cmds := make([]Command, 0, len(commandRegistry))
	for _, cmd := range commandRegistry {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func init() {
	// 내장 명령어 등록
	Register(&PingCommand{})
	Register(&HelpCommand{})
}
