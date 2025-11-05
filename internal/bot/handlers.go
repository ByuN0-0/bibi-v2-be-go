package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"bibi-bot-v2/internal/commands"
)

// registerHandlers 모든 이벤트 핸들러를 등록합니다
func (b *Bot) registerHandlers() {
	b.Session.AddHandler(b.onReady)
	b.Session.AddHandler(b.onMessageCreate)
}

// onReady Bot이 준비 완료되었을 때 호출됩니다
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	b.Logger.Info(fmt.Sprintf("Bot 준비 완료. 로그인: %s", event.User.Username))
	s.UpdateGameStatus(0, "!help")
}

// onMessageCreate 메시지가 생성되었을 때 호출됩니다
func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot이 보낸 메시지 무시
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 명령어 처리
	if strings.HasPrefix(m.Content, "!") {
		b.handleCommand(s, m)
	}
}

// handleCommand 명령어를 처리합니다
func (b *Bot) handleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) == 0 {
		return
	}

	cmdName := strings.TrimPrefix(args[0], "!")
	cmdArgs := args[1:]

	// 등록된 명령어 실행
	cmd := commands.GetCommand(cmdName)
	if cmd == nil {
		b.Logger.Info(fmt.Sprintf("알 수 없는 명령어: %s", cmdName))
		return
	}

	err := cmd.Execute(s, m, cmdArgs)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("명령어 실행 실패: %v", err))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("오류: %v", err))
	}
}
