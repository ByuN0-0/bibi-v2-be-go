package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"bibi-bot-v2/internal/commands"
)

// initializeCommands 명령어들을 초기화합니다
func (b *Bot) initializeCommands() {
	// Weather 명령어 등록
	commands.Register(commands.NewWeatherCommand(b.Config.WeatherAPIKey))
}

// registerHandlers 모든 이벤트 핸들러를 등록합니다
func (b *Bot) registerHandlers() {
	b.Session.AddHandler(b.onReady)
	b.Session.AddHandler(b.onInteractionCreate)
}

// onReady Bot이 준비 완료되었을 때 호출됩니다
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	b.Logger.Info(fmt.Sprintf("Bot 준비 완료. 로그인: %s", event.User.Username))
	startTime := time.Now().Format("2006-01-02T15:04")
	s.UpdateGameStatus(0, startTime)
}

// onInteractionCreate 슬래쉬 명령어가 실행되었을 때 호출됩니다
func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ApplicationCommand 타입만 처리
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	cmdName := i.ApplicationCommandData().Name

	// 등록된 명령어 실행
	cmd := commands.GetCommand(cmdName)
	if cmd == nil {
		b.Logger.Info(fmt.Sprintf("알 수 없는 명령어: %s", cmdName))
		return
	}

	err := cmd.Execute(s, i)
	if err != nil {
		b.Logger.Error(fmt.Sprintf("명령어 실행 실패: %v", err))
		// 에러 응답
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("오류: %v", err),
			},
		})
	}
}
