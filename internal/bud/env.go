package bud

type Env map[string]string

func (e Env) Add(env Env) Env {
	out := Env{}
	for k, v := range e {
		out[k] = v
	}
	for k, v := range env {
		out[k] = v
	}
	return out
}

func (e Env) List() (env []string) {
	for k, v := range e {
		env = append(env, k+"="+v)
	}
	return env
}
