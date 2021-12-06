package cli

import (
	"context"
	"net/http"
)

type Serve struct {
	Port int
}

func (s *Serve) Run(ctx context.Context, web http.Handler) error {
	return nil
}
