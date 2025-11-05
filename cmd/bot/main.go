package main

import (
	"fmt"
	"os"

	"bibi-bot-v2/internal/bot"
	"bibi-bot-v2/internal/config"
	"bibi-bot-v2/internal/logger"
)

func main() {
	// 로거 초기화
	log := logger.NewLogger()
	log.Info("bibi-bot-v2 시작 중...")

	// 설정 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(fmt.Sprintf("설정 로드 실패: %v", err))
		os.Exit(1)
	}

	// Bot 초기화
	botInstance, err := bot.NewBot(cfg, log)
	if err != nil {
		log.Error(fmt.Sprintf("Bot 초기화 실패: %v", err))
		os.Exit(1)
	}

	// Bot 시작
	if err := botInstance.Start(); err != nil {
		log.Error(fmt.Sprintf("Bot 시작 실패: %v", err))
		os.Exit(1)
	}
}
