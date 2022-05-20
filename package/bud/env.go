package bud

import "sort"

type Env map[string]string

func (e Env) List() (env []string) {
	for k, v := range e {
		env = append(env, k+"="+v)
	}
	sort.Strings(env)
	return env
}
