package cli

import "sort"

// Env is an environment helper type
type Env map[string]string

// List out the environment keys for use with cmd.Env in alphanumerical order.
func (e Env) List() (env []string) {
	for k, v := range e {
		env = append(env, k+"="+v)
	}
	sort.Strings(env)
	return env
}

// Defaults sets `e[key]` to `env[key]` if `e[key]` doesn't exist.
func (e Env) Defaults(env Env) Env {
	if e == nil {
		e = Env{}
	}
	for k, v := range env {
		if _, ok := e[k]; ok {
			continue
		}
		e[k] = v
	}
	return e
}
