package config

import (
	"fmt"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
)

type Config struct {
	Token         string `json:"token"`
	ClientID      string `json:"clientId"`
	Secret        string `json:"secret"`
	PublicKey     string `json:"publicKey"`
	WeatherAPIKey string `json:"weatherApiKey"`
}

// Validate 설정값을 검증합니다
func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Token, validation.Required.Error("DISCORD_BOT_TOKEN은 필수입니다")),
		validation.Field(&c.ClientID, validation.Required.Error("DISCORD_CLIENT_ID는 필수입니다")),
		validation.Field(&c.Secret, validation.Required.Error("DISCORD_CLIENT_SECRET은 필수입니다")),
		validation.Field(&c.PublicKey, validation.Required.Error("DISCORD_PUBLIC_KEY는 필수입니다")),
		validation.Field(&c.WeatherAPIKey, validation.Required.Error("OPENWEATHER_API_KEY는 필수입니다")),
	)
}

// LoadConfig 설정을 .env 파일에서 로드합니다
func LoadConfig() (*Config, error) {
	// .env 파일 로드 (실패해도 무시)
	_ = godotenv.Load()

	config := &Config{
		Token:         os.Getenv("DISCORD_BOT_TOKEN"),
		ClientID:      os.Getenv("DISCORD_CLIENT_ID"),
		Secret:        os.Getenv("DISCORD_CLIENT_SECRET"),
		PublicKey:     os.Getenv("DISCORD_PUBLIC_KEY"),
		WeatherAPIKey: os.Getenv("OPENWEATHER_API_KEY"),
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("설정 검증 실패: %w", err)
	}

	return config, nil
}
