package di_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/livebud/bud/pkg/di"
	"github.com/matryer/is"
)

type Env struct {
	name string
}

func loadEnv(in di.Injector) (*Env, error) {
	return &Env{"production"}, nil
}

type Log struct {
	env *Env
	lvl string
}

func loadLog(in di.Injector) (*Log, error) {
	env, err := di.Load[*Env](in)
	if err != nil {
		return nil, err
	}
	return &Log{env: env, lvl: "info"}, nil
}

func TestDI(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader(in, loadEnv)
	di.Loader(in, loadLog)
	log, err := di.Load[*Log](in)
	is.NoErr(err)
	is.Equal(log.lvl, "info")
}

func TestClone(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader(in, loadEnv)
	in2 := di.Clone(in)
	env, err := di.Load[*Env](in2)
	is.NoErr(err)
	is.Equal(env.name, "production")
	di.Loader(in2, loadLog)
	log, err := di.Load[*Log](in2)
	is.NoErr(err)
	is.Equal(log.lvl, "info")
	log, err = di.Load[*Log](in)
	is.True(err != nil)
	is.True(errors.Is(err, di.ErrNoLoader))
	is.Equal(log, nil)
}

func ExampleLoader() {
	in := di.New()
	di.Loader(in, loadEnv)
	di.Loader(in, loadLog)
	log, err := di.Load[*Log](in)
	if err != nil {
		panic(err)
	}
	fmt.Println(log.env.name)
	// Output: production
}
