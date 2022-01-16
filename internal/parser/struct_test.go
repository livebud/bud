package parser_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/modtest"
	"gitlab.com/mnm/bud/internal/parser"
)

func TestStructMethod(t *testing.T) {
	is := is.New(t)
	module := modtest.Make(t, modtest.Module{
		Files: map[string][]byte{
			"go.mod": []byte(`module app.com/app`),
			"app.go": []byte(`
				package app

				type A struct {
				}

				func (a *A) Method() {}
			`),
		},
	})
	p := parser.New(module)
	pkg, err := p.Parse(".")
	is.NoErr(err)
	is.Equal(pkg.Name(), "app")
	stct := pkg.Struct("A")
	is.True(stct != nil)
	method := stct.Method("method")
	is.True(method == nil)
	method = stct.Method("methods")
	is.True(method == nil)
	method = stct.Method("Method")
	is.True(method != nil)
}

func TestStructTag(t *testing.T) {
	is := is.New(t)
	module := modtest.Make(t, modtest.Module{
		Files: map[string][]byte{
			"go.mod": []byte(`module app.com/app`),
			"app.go": []byte(`
				package app

				type User struct {
					Some string ` + "`" + `json:"some,omitempty" is:"email"` + "`" + `
					None bool
				}
			`),
		},
	})
	p := parser.New(module)
	pkg, err := p.Parse(".")
	is.NoErr(err)
	is.Equal(pkg.Name(), "app")
	stct := pkg.Struct("User")
	is.True(stct != nil)
	// User.Some field
	field := stct.Field("Some")
	is.True(field != nil)
	tags, err := field.Tags()
	is.NoErr(err)
	is.Equal(len(tags), 2)
	is.Equal(tags[0].Key, "json")
	is.Equal(tags[0].Value, "some")
	is.Equal(len(tags[0].Options), 1)
	is.Equal(tags[0].Options[0], "omitempty")
	is.Equal(tags[1].Key, "is")
	is.Equal(tags[1].Value, "email")
	is.Equal(len(tags[1].Options), 0)
	// User.None field
	field = stct.Field("None")
	is.True(field != nil)
	tags, err = field.Tags()
	is.NoErr(err)
	is.Equal(len(tags), 0)
}
