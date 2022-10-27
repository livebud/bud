package log

import (
	"fmt"
)

// Level of the logger
type Level uint8

// Log level
const (
	DebugLevel Level = iota + 1
	InfoLevel
	NoticeLevel
	WarnLevel
	ErrorLevel
)

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case NoticeLevel:
		return "notice"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return ""
	}
}

func ParseLevel(level string) (Level, error) {
	switch level {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "notice":
		return NoticeLevel, nil
	case "warn":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	}
	return 0, fmt.Errorf("log: %q is not a valid level", level)
}
