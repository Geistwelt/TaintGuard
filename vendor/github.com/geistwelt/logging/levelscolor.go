package logging

import (
	"fmt"
)

///////////////////////////////////////////////////////////////////
// LogLevel 日志记录级别。

type LogLevel int8

const (
	PanicLevel = iota + 1
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func (l LogLevel) CapitalString() string {
	switch l {
	case PanicLevel:
		return "PANIC"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN "
	case InfoLevel:
		return "INFO "
	case DebugLevel:
		return "DEBUG"
	default:
		panic(fmt.Sprintf("invalid log level: (%d)", l))
	}
}

func (l LogLevel) LowercaseString() string {
	switch l {
	case PanicLevel:
		return "panic"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn "
	case InfoLevel:
		return "info "
	case DebugLevel:
		return "debug"
	default:
		panic(fmt.Sprintf("invalid log level: (%d)", l))
	}
}

func (l LogLevel) SpecifiedColor() LevelColor {
	switch l {
	case DebugLevel:
		return DebugLevelColor
	case InfoLevel:
		return InfoLevelColor
	case WarnLevel:
		return WarnLevelColor
	case ErrorLevel:
		return ErrorLevelColor
	case PanicLevel:
		return PanicLevelColor
	default:
		panic(fmt.Sprintf("invalid log level: (%d)", l))
	}
}

///////////////////////////////////////////////////////////////////
// LevelColor 为各个日志记录级别定义的颜色。

type LevelColor int8

const (
	DebugLevelColor LevelColor = 34 // 蓝色
	InfoLevelColor  LevelColor = 32 // 绿色
	WarnLevelColor  LevelColor = 33 // 黄色
	ErrorLevelColor LevelColor = 31 // 红色
	PanicLevelColor LevelColor = 35 // 紫色
)

func (lc LevelColor) Color() string {
	return fmt.Sprintf("\x1b[%dm", lc)
}

func ResetColor() string {
	return "\x1b[0m"
}
