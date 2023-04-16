// Package errs makes it easier to join multiple errors into a single error.
package errs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/livebud/bud/internal/ansi"
)

// Join multiple errors together into one error
func Join(errs ...error) error {
	var agg error
	for _, err := range errs {
		if err == nil {
			continue
		} else if agg == nil {
			agg = err
			continue
		} else if errors.Is(err, agg) {
			agg = fmt.Errorf("%w. %s", agg, err)
		} else {
			agg = fmt.Errorf("%s. %s", agg, err)
		}
	}
	return agg
}

// Errors is an optional interface that be used to unwrap multiple errors
type Errors interface {
	Errors() []error
}

// Format reverses the error order to make the cause come first
func Format(err error) string {
	// Most errors in Bud are joined by a period
	lines := strings.Split(err.Error(), ". ")
	lineLen := len(lines)
	stack := make([]string, lineLen)
	j := lineLen - 1
	// Reverse the error order
	for i := 0; i < lineLen; i++ {
		line := lines[j]
		if i > 0 {
			line = "  " + ansi.Dim(line)
		}
		stack[i] = line
		j--
	}
	return strings.Join(stack, "\n")
}
