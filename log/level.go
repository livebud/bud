package log

import (
	"fmt"
)

// Level of the logger
type Level uint8

// Log level
const (
	LevelDebug Level = iota + 1
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
)

func (level Level) String() string {
	switch level {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelNotice:
		return "notice"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return ""
	}
}

func ParseLevel(level string) (Level, error) {
	switch level {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "notice":
		return LevelNotice, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	}
	return 0, fmt.Errorf("log: %q is not a valid level", level)
}
