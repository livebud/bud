// Package errs makes it easier to join multiple errors into a single error.
package errs

import (
	"errors"
	"fmt"
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
