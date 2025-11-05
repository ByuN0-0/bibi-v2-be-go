package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token     string
	ClientID  string
	Secret    string
	PublicKey string
}

// LoadConfig 설정을 .env 파일에서 로드합니다
func LoadConfig() (*Config, error) {
	// .env 파일 로드 (실패해도 무시)
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN 환경변수가 설정되지 않음")
	}

	return &Config{
		Token:     token,
		ClientID:  os.Getenv("DISCORD_CLIENT_ID"),
		Secret:    os.Getenv("DISCORD_CLIENT_SECRET"),
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
	}, nil
}
