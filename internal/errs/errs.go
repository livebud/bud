// Package errs makes it easier to join multiple errors into a single error.
package errs

import "fmt"

// Join multiple errors together into one error
func Join(errors ...error) error {
	var agg error
	for _, err := range errors {
		if err == nil {
			continue
		} else if agg == nil {
			agg = err
			continue
		} else {
			agg = fmt.Errorf("%s. %s", agg, err)
		}
	}
	return agg
}
