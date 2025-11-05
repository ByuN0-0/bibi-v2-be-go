package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

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

	// 무한 대기
	select {}
}

// Stop Bot을 종료합니다
func (b *Bot) Stop() error {
	b.Logger.Info("Bot 종료 중...")
	return b.Session.Close()
}
