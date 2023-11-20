package env

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

// Load an environment from a struct
func Load[Env any](e Env) (Env, error) {
	if err := env.Parse(e); err != nil {
		return e, fmt.Errorf("env: unable to load the environment: %w", err)
	}
	return e, nil
}
