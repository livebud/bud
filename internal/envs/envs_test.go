package envs_test

import (
	"testing"

	"github.com/livebud/bud/internal/envs"
	"github.com/livebud/bud/internal/is"
)

func TestFrom(t *testing.T) {
	is := is.New(t)
	env := envs.From([]string{
		"HOME=/users/matt",
		"PATH=/usr/local/bin:/usr/bin",
	})
	is.Equal(env["HOME"], "/users/matt")
	is.Equal(env["PATH"], "/usr/local/bin:/usr/bin")
}

func TestList(t *testing.T) {
	is := is.New(t)
	env := envs.Map{
		"B": "B",
		"A": "A",
		"C": "C",
	}
	list := env.List()
	is.Equal(len(list), 3)
	is.Equal(list[0], "A=A")
	is.Equal(list[1], "B=B")
	is.Equal(list[2], "C=C")
}
