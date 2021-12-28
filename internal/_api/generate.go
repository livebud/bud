package api

import (
	"context"

	"gitlab.com/mnm/bud/fsync"
)

// Generate the code
func (a *API) Generate(ctx context.Context) error {
	return fsync.Dir(a.genFS, ".", a.appFS, ".")
}
