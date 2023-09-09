package env

// Bud environment variables
type Bud struct {
	Env    string `env:"BUD_ENV" default:"development"`
	Listen string `env:"BUD_LISTEN" default:":3000"`
}
