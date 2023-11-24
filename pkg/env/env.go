package env

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

// Load an environment from a struct
func Load[Env any](e Env) (Env, error) {
	opts := env.Options{
		RequiredIfNoDef: true,
	}
	if err := env.ParseWithOptions(e, opts); err != nil {
		return e, fmt.Errorf("env: unable to load the environment: %w", err)
	}
	return e, nil
}
