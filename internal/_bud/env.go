package bud

type Env map[string]string

func (e Env) List() (env []string) {
	for k, v := range e {
		env = append(env, k+"="+v)
	}
	return env
}
