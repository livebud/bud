package envs

import (
	"sort"
	"strings"
)

// Map is an environment helper type
type Map map[string]string

// Parse an environment list into a map
func From(list []string) Map {
	m := Map{}
	for _, row := range list {
		kvs := strings.SplitN(row, "=", 2)
		if len(kvs) != 2 {
			continue
		}
		m[kvs[0]] = kvs[1]
	}
	return m
}

// List out the environment keys in alphanumerical order.
func (m Map) List() (env []string) {
	for k, v := range m {
		env = append(env, k+"="+v)
	}
	sort.Strings(env)
	return env
}

func (m Map) Append(list ...string) Map {
	lm := From(list)
	for k, v := range lm {
		m[k] = v
	}
	return m
}
