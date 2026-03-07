package utils

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	LogTrace LogLevel = iota
	LogDebug
	LogProc
	LogInfo
	LogWarn
	LogError
)

var LogLevelStr = map[LogLevel]string{
	LogTrace: "\033[35mTRACE\033[0m",
	LogDebug: "\033[36mDEBUG\033[0m",
	LogInfo:  "\033[32mINFO\033[0m",
	LogProc:  "\033[34mPROCESS\033[0m",
	LogWarn:  "\033[33mWARN\033[0m",
	LogError: "\033[31mERROR\033[0m",
}

type Logger struct {
	out          *os.File
	minimumLevel LogLevel
}

func NewLogger(out *os.File) *Logger {
	return &Logger{out: out}
}

var GlobalLogger *Logger = NewLogger(os.Stderr)

func (l *Logger) Log(msg string, level LogLevel) {
	if level > l.minimumLevel {
		fmt.Fprintf(l.out, "\033[90m%s\033[0m %s: %s\n", time.Now().Format(time.DateTime), LogLevelStr[level], msg)
	}
}

func (l *Logger) MinimumLevel(min LogLevel) {
	l.minimumLevel = min
}

func (l *Logger) SetNewOutput(newOut *os.File) {
	l.out = newOut
}
