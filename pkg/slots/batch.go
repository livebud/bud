package slots

import "net/http"

// Batch a list of handlers to be called in parallel.
func Batch(handlers ...http.Handler) http.Handler {
	// TODO: finish me
	return handlers[0]
}
