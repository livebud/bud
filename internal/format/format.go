package format

import (
	"fmt"
	"strings"

	"github.com/lithammer/dedent"
)

func format(message string) string {
	return strings.TrimSpace(dedent.Dedent(message)) + "\n"
}

func Sprintf(message string, args ...interface{}) string {
	return fmt.Sprintf(format(message), args...)
}

func Errorf(message string, args ...interface{}) error {
	return fmt.Errorf(format(message), args...)
}
