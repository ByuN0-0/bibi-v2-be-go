package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// OneCallResponse OpenWeather One Call API 3.0 응답 구조
type OneCallResponse struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`
	Current  Current `json:"current"`
	Daily    []Daily `json:"daily"`
	Alerts   []Alert `json:"alerts"`
}

// Current 현재 날씨 정보
type Current struct {
	DT        int64       `json:"dt"`
	Temp      float64     `json:"temp"`
	FeelsLike float64     `json:"feels_like"`
	Pressure  int         `json:"pressure"`
	Humidity  int         `json:"humidity"`
	Clouds    int         `json:"clouds"`
	Visibility int        `json:"visibility"`
	WindSpeed float64     `json:"wind_speed"`
	WindDeg   int         `json:"wind_deg"`
	WindGust  float64     `json:"wind_gust"`
	Weather   []Weather   `json:"weather"`
	Sunrise   int64       `json:"sunrise"`
	Sunset    int64       `json:"sunset"`
	UVI       float64     `json:"uvi"`
}

// Daily 일일 날씨 정보
type Daily struct {
	DT        int64       `json:"dt"`
	Temp      TempInfo    `json:"temp"`
	FeelsLike FeelsLike   `json:"feels_like"`
	Pressure  int         `json:"pressure"`
	Humidity  int         `json:"humidity"`
	Clouds    int         `json:"clouds"`
	WindSpeed float64     `json:"wind_speed"`
	WindDeg   int         `json:"wind_deg"`
	WindGust  float64     `json:"wind_gust"`
	Weather   []Weather   `json:"weather"`
	Pop       float64     `json:"pop"`
	Rain      float64     `json:"rain"`
	Snow      float64     `json:"snow"`
	UVI       float64     `json:"uvi"`
	Summary   string      `json:"summary"`
}

// TempInfo 온도 정보
type TempInfo struct {
	Day   float64 `json:"day"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

// FeelsLike 체감 온도
type FeelsLike struct {
	Day   float64 `json:"day"`
	Night float64 `json:"night"`
	Eve   float64 `json:"eve"`
	Morn  float64 `json:"morn"`
}

// Weather 날씨 상태
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Alert 기상 경보
type Alert struct {
	SenderName  string `json:"sender_name"`
	Event       string `json:"event"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Description string `json:"description"`
}

// WeatherClient OpenWeather API 클라이언트
type WeatherClient struct {
	apiKey string
	client *http.Client
}

// NewWeatherClient WeatherClient를 생성합니다
func NewWeatherClient(apiKey string) *WeatherClient {
	return &WeatherClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// GetSeoulWeather 서울의 현재 날씨와 일일 예보를 조회합니다
func (wc *WeatherClient) GetSeoulWeather() (*OneCallResponse, error) {
	// 서울 좌표
	lat := 37.5665
	lon := 126.9780

	baseURL := "https://api.openweathermap.org/data/3.0/onecall"
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lon))
	params.Add("appid", wc.apiKey)
	params.Add("units", "metric")
	params.Add("lang", "ko")
	params.Add("exclude", "minutely,hourly")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := wc.client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("날씨 조회 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 오류 (상태코드: %d): %s", resp.StatusCode, string(body))
	}

	var weather OneCallResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	return &weather, nil
}
