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

func loadEnv() (*Env, error) {
	return &Env{"production"}, nil
}

type Log struct {
	env *Env
	lvl string
}

func (l *Log) Info(msg string) {
}

func loadLog(env *Env) (*Log, error) {
	return &Log{env: env, lvl: "info"}, nil
}

type Logger interface {
	Info(msg string)
}

type DB struct {
	env *Env
	log Logger
}

func loadDB(env *Env, log Logger) (*DB, error) {
	return &DB{env: env, log: log}, nil
}

func TestDI(t *testing.T) {
	is := is.New(t)
	in := di.New()
	is.NoErr(di.Provide[*Env](in, loadEnv))
	is.NoErr(di.Provide[*Log](in, loadLog))
	log, err := di.Load[*Log](in)
	is.NoErr(err)
	is.Equal(log.lvl, "info")
}

type Router struct {
	log    *Log
	routes map[string]func()
}

func (r *Router) Get(path string, fn func()) {
	r.routes[path] = fn
}

func loadRouter(log *Log) (*Router, error) {
	return &Router{log: log, routes: map[string]func(){}}, nil
}

func loadUsers() *Users {
	return &Users{}
}

type Users struct{}

func (*Users) Index() {}

func userRoutes(users *Users, router *Router) error {
	router.Get("/users", users.Index)
	return nil
}

func loadPosts() *Posts {
	return &Posts{}
}

type Posts struct{}

func (*Posts) Show() {}

func postRoutes(router *Router, posts *Posts) {
	router.Get("/posts", posts.Show)
}

func TestWhen(t *testing.T) {
	is := is.New(t)
	in := di.New()
	is.NoErr(di.Provide[*Env](in, loadEnv))
	is.NoErr(di.Provide[*Log](in, loadLog))
	is.NoErr(di.Provide[*Router](in, loadRouter))
	is.NoErr(di.Provide[*Users](in, loadUsers))
	is.NoErr(di.Provide[*Posts](in, loadPosts))
	is.NoErr(di.When[*Router](in, userRoutes))
	is.NoErr(di.When[*Router](in, postRoutes))
	router, err := di.Load[*Router](in)
	is.NoErr(err)
	is.Equal(router.log.lvl, "info")
	is.Equal(len(router.routes), 2)
}

func TestCache(t *testing.T) {
	is := is.New(t)
	in := di.New()
	is.NoErr(di.Provide[*Env](in, loadEnv))
	is.NoErr(di.Provide[*Log](in, loadLog))
	log1, err := di.Load[*Log](in)
	is.NoErr(err)
	log2, err := di.Load[*Log](in)
	is.NoErr(err)
	is.Equal(log1, log2)
}

func TestClone(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Provide[*Env](in, loadEnv)
	in2 := di.Clone(in)
	env, err := di.Load[*Env](in2)
	is.NoErr(err)
	is.Equal(env.name, "production")
	di.Provide[*Log](in2, loadLog)
	log, err := di.Load[*Log](in2)
	is.NoErr(err)
	is.Equal(log.lvl, "info")
	log, err = di.Load[*Log](in)
	is.True(err != nil)
	is.True(errors.Is(err, di.ErrNoProvider))
	is.Equal(log, nil)
}

func TestInterface(t *testing.T) {
	is := is.New(t)
	in := di.New()
	is.NoErr(di.Provide[*Env](in, loadEnv))
	is.NoErr(di.Provide[Logger](in, loadLog))
	is.NoErr(di.Provide[*DB](in, loadDB))
	log, err := di.Load[Logger](in)
	is.NoErr(err)
	is.True(log != nil)
	db, err := di.Load[*DB](in)
	is.NoErr(err)
	is.True(db != nil)
}

func ExampleLoad() {
	in := di.New()
	di.Provide[*Env](in, loadEnv)
	di.Provide[*Log](in, loadLog)
	log, err := di.Load[*Log](in)
	if err != nil {
		panic(err)
	}
	fmt.Println(log.env.name)
	// Output: production
}
