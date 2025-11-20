package commands

import (
	"fmt"
	"time"

	"bibi-bot-v2/internal/services"
	"github.com/bwmarrin/discordgo"
)

type WeatherCommand struct {
	weatherClient *services.WeatherClient
}

// NewWeatherCommand WeatherCommand를 생성합니다
func NewWeatherCommand(apiKey string) *WeatherCommand {
	return &WeatherCommand{
		weatherClient: services.NewWeatherClient(apiKey),
	}
}

func (wc *WeatherCommand) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "오늘날씨",
		Description: "서울의 현재 날씨를 조회합니다",
	}
}

func (wc *WeatherCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// 날씨 데이터 조회
	weather, err := wc.weatherClient.GetSeoulWeather()
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("⚠️ 날씨 조회 실패: %v", err),
			},
		})
	}

	// 날씨 데이터 검증
	if len(weather.Current.Weather) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ 날씨 데이터를 가져올 수 없습니다",
			},
		})
	}

	// 현재 날씨 Embed 생성
	currentEmbed := &discordgo.MessageEmbed{
		Title:       "🌤️ 서울의 현재 날씨",
		Description: weather.Current.Weather[0].Description,
		Color:       0x3498DB,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🌡️ 온도",
				Value:  fmt.Sprintf("%.1f°C (체감: %.1f°C)", weather.Current.Temp, weather.Current.FeelsLike),
				Inline: true,
			},
			{
				Name:   "💧 습도",
				Value:  fmt.Sprintf("%d%%", weather.Current.Humidity),
				Inline: true,
			},
			{
				Name:   "💨 풍속",
				Value:  fmt.Sprintf("%.1f m/s", weather.Current.WindSpeed),
				Inline: true,
			},
			{
				Name:   "🌬️ 돌풍",
				Value:  fmt.Sprintf("%.1f m/s", weather.Current.WindGust),
				Inline: true,
			},
			{
				Name:   "🌥️ 기압",
				Value:  fmt.Sprintf("%d hPa", weather.Current.Pressure),
				Inline: true,
			},
			{
				Name:   "☁️ 구름",
				Value:  fmt.Sprintf("%d%%", weather.Current.Clouds),
				Inline: true,
			},
			{
				Name:   "🌞 자외선지수",
				Value:  fmt.Sprintf("%.1f", weather.Current.UVI),
				Inline: true,
			},
			{
				Name:   "🔍 시정",
				Value:  fmt.Sprintf("%d m", weather.Current.Visibility),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 내일 날씨 Embed 생성 (일일 예보 첫 번째 항목)
	var embeds []*discordgo.MessageEmbed
	embeds = append(embeds, currentEmbed)

	if len(weather.Daily) > 1 {
		tomorrow := weather.Daily[1]
		tomorrowEmbed := &discordgo.MessageEmbed{
			Title:       "📅 내일 날씨 예보",
			Description: tomorrow.Summary,
			Color:       0x2ECC71,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "🌡️ 온도",
					Value:  fmt.Sprintf("낮: %.1f°C | 밤: %.1f°C\n최고: %.1f°C | 최저: %.1f°C", tomorrow.Temp.Day, tomorrow.Temp.Night, tomorrow.Temp.Max, tomorrow.Temp.Min),
					Inline: false,
				},
				{
					Name:   "💧 습도",
					Value:  fmt.Sprintf("%d%%", tomorrow.Humidity),
					Inline: true,
				},
				{
					Name:   "💨 풍속",
					Value:  fmt.Sprintf("%.1f m/s", tomorrow.WindSpeed),
					Inline: true,
				},
				{
					Name:   "☔ 강수확률",
					Value:  fmt.Sprintf("%.0f%%", tomorrow.Pop*100),
					Inline: true,
				},
				{
					Name:   "🌞 자외선지수",
					Value:  fmt.Sprintf("%.1f", tomorrow.UVI),
					Inline: true,
				},
			},
		}
		embeds = append(embeds, tomorrowEmbed)
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})
}
