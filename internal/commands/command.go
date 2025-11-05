package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Command 명령어 인터페이스
type Command interface {
	// Name 명령어 이름을 반환합니다
	Name() string

	// Description 명령어 설명을 반환합니다
	Description() string

	// Execute 명령어를 실행합니다
	Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error
}

// CommandRegistry 등록된 명령어들
var commandRegistry = make(map[string]Command)

// Register 명령어를 등록합니다
func Register(cmd Command) {
	commandRegistry[cmd.Name()] = cmd
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
