# Bibi Bot v2

Discord bot을 Go로 다시 구현한 프로젝트입니다. (v1은 Java)

## 개요

Bibi Bot v2는 Go로 구현된 고성능 Discord 봇입니다.

## 스택

- **Language**: Go 1.20+
- **Discord Library**: discordgo
- **Infrastructure**: Terraform
- **Cloud**: Google Cloud Platform (GCP)
- **Compute**: Compute Engine (e2-micro, us-west1)
- **Registry**: Google Container Registry (GCR)
- **Container OS**: Container-Optimized OS (COS)

## 기능

- 기본 봇 커맨드 처리
- 이벤트 기반 핸들러
- 모듈식 구조로 쉬운 커맨드 추가

## 프로젝트 구조 (예시)

```
bibi-bot-v2/
├── cmd/
│   └── bot/
│       └── main.go                # 진입점
├── internal/
│   ├── bot/
│   │   ├── bot.go                 # Bot 구조체 및 관리
│   │   └── handlers.go            # 이벤트 핸들러
│   ├── commands/
│   │   ├── command.go             # 커맨드 인터페이스
│   │   ├── ping.go
│   │   └── help.go
│   ├── config/
│   │   └── config.go              # 설정 로딩
│   └── logger/
│       └── logger.go              # 로깅
├── go.mod
├── go.sum
└── README.md
```

## 시작하기

### 요구사항

- Go 1.20 이상
- Discord Bot Token
- GCP 계정 (배포 시)
- Docker (배포 시)
- Terraform (배포 시)

### 로컬 실행

1. 의존성 설치

```bash
go mod download
```

2. 환경 변수 설정

```bash
export DISCORD_TOKEN="your-bot-token"
```

3. 실행

```bash
go run cmd/bot/main.go
```

## 환경 변수

| 변수            | 설명             | 필수 |
| --------------- | ---------------- | ---- |
| `DISCORD_TOKEN` | Discord Bot 토큰 | ✓    |

## 커맨드 추가하기

1. `internal/commands/` 디렉토리에 새 파일 생성

```go
// internal/commands/mycommand.go
package commands

import "github.com/bwmarrin/discordgo"

type MyCommand struct{}

func (c *MyCommand) Name() string {
	return "mycommand"
}

func (c *MyCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Response")
}
```

2. Bot 핸들러에 등록

```go
// internal/bot/handlers.go
commands := []commands.Command{
	&commands.MyCommand{},
}
```

## 라이선스

MIT
