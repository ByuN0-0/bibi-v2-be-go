package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"bibi-bot-v2/internal/commands"
	"bibi-bot-v2/internal/config"
	"bibi-bot-v2/internal/logger"
)

type Bot struct {
	Session *discordgo.Session
	Config  *config.Config
	Logger  *logger.Logger
}

// NewBot Bot 구조체를 생성합니다
func NewBot(cfg *config.Config, log *logger.Logger) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("Discord 세션 생성 실패: %w", err)
	}

	bot := &Bot{
		Session: session,
		Config:  cfg,
		Logger:  log,
	}

	// 명령어 초기화
	bot.initializeCommands()

	// 핸들러 등록
	bot.registerHandlers()

	return bot, nil
}

// Start Bot을 시작합니다
func (b *Bot) Start() error {
	b.Logger.Info("Discord에 연결 중...")

	// Intent 설정
	b.Session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("Discord 연결 실패: %w", err)
	}

	b.Logger.Info("Bot이 성공적으로 연결되었습니다")

	// 슬래쉬 명령어 등록
	err = b.registerSlashCommands()
	if err != nil {
		return fmt.Errorf("슬래쉬 명령어 등록 실패: %w", err)
	}

	// 무한 대기
	select {}
}

// registerSlashCommands 슬래쉬 명령어를 Discord에 등록합니다
func (b *Bot) registerSlashCommands() error {
	b.Logger.Info("슬래쉬 명령어 등록 중...")

	for _, cmd := range commands.GetAllCommands() {
		appCmd := cmd.ApplicationCommand()
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", appCmd)
		if err != nil {
			return fmt.Errorf("명령어 '%s' 등록 실패: %w", appCmd.Name, err)
		}
		b.Logger.Info(fmt.Sprintf("명령어 등록 완료: /%s", appCmd.Name))
	}

	b.Logger.Info("모든 슬래쉬 명령어 등록 완료")
	return nil
}

// Stop Bot을 종료합니다
func (b *Bot) Stop() error {
	b.Logger.Info("Bot 종료 중...")
	return b.Session.Close()
}
