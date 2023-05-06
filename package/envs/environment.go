package envs

import environment "github.com/caarlos0/env/v8"

func Load[Env any]() (env *Env, err error) {
	env = new(Env)
	if err := environment.Parse(env); err != nil {
		return nil, err
	}
	return env, nil
}
