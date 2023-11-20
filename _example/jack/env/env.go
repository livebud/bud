package env

import (
	"github.com/livebud/bud/pkg/env"
)

func Load() (*Env, error) {
	return env.Load(new(Env))
}

type Env struct {
	API_URL            string `env:"API_URL,required"`
	SLACK_CLIENT_ID    string `env:"SLACK_CLIENT_ID,required"`
	SLACK_REDIRECT_URL string `env:"SLACK_REDIRECT_URL,required"`
	SLACK_SCOPE        string `env:"SLACK_SCOPE,required"`
	SLACK_USER_SCOPE   string `env:"SLACK_USER_SCOPE,required"`
	STRIPE_CLIENT_KEY  string `env:"STRIPE_CLIENT_KEY,required"`
}
