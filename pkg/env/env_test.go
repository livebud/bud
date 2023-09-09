package env_test

import (
	"strings"
	"testing"

	"github.com/livebud/bud/pkg/env"
	"github.com/matryer/is"
)

func TestInvalid(t *testing.T) {
	is := is.New(t)
	s, err := env.Load[string]()
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "expected a pointer to a struct"))
	is.Equal(s, "")
}

type Env struct {
	Database Database
	Debug    bool `env:"DEBUG"`
	Port     int  `env:"PORT" default:"3000"`
}

type Database struct {
	URL string `env:"DATABASE_URL" default:"postgres://localhost:5432/blog?sslmode=disable"`
}

func (d *Database) Development() error {
	d.URL = strings.ReplaceAll(d.URL, "blog", "dev")
	return nil
}

func (e *Env) Development() error {
	e.Port = e.Port + 3000
	return nil
}

func TestEnv(t *testing.T) {
	t.Skip("Make reflection less fragile")
	t.Setenv("BUD_ENV", "development")
	t.Setenv("DEBUG", "true")
	is := is.New(t)
	e, err := env.Load[Env]()
	is.NoErr(err)
	is.Equal(e.Port, 6000)
	is.Equal(e.Database.URL, "postgres://localhost:5432/dev?sslmode=disable")
	is.Equal(e.Debug, true)
}

func TestEnvPtr(t *testing.T) {
	t.Setenv("BUD_ENV", "development")
	t.Setenv("DEBUG", "true")
	is := is.New(t)
	e, err := env.Load[*Env]()
	is.NoErr(err)
	is.Equal(e.Port, 6000)
	is.Equal(e.Database.URL, "postgres://localhost:5432/dev?sslmode=disable")
	is.Equal(e.Debug, true)
}

func TestMissing(t *testing.T) {
	type Env struct {
		Port int `env:"PORT"`
	}
	is := is.New(t)
	e, err := env.Load[*Env]()
	is.True(err != nil)
	is.Equal(err.Error(), `env: missing required environment variable "PORT"`)
	is.Equal(e, nil)
}
