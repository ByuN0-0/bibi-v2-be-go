package logger

import (
	"fmt"
	"time"
)

// Logger 로깅을 관리합니다
type Logger struct {
	debug bool
}

// NewLogger Logger를 생성합니다
func NewLogger() *Logger {
	return &Logger{
		debug: false,
	}
}

// Info 정보 레벨 로그를 출력합니다
func (l *Logger) Info(message string) {
	l.log("INFO", message)
}

// Error 에러 레벨 로그를 출력합니다
func (l *Logger) Error(message string) {
	l.log("ERROR", message)
}

// Debug 디버그 레벨 로그를 출력합니다
func (l *Logger) Debug(message string) {
	if l.debug {
		l.log("DEBUG", message)
	}
}

// Warn 경고 레벨 로그를 출력합니다
func (l *Logger) Warn(message string) {
	l.log("WARN", message)
}

// log 로그를 내부적으로 출력합니다
func (l *Logger) log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s: %s\n", timestamp, level, message)
}

// SetDebug 디버그 모드를 설정합니다
func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}
