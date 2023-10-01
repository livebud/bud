package env

type Env struct {
	Port int    `env:"PORT" default:"8080"`
	Log  string `env:"LOG" default:"info"`
}
