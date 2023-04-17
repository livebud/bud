package env

type Env struct {
	Database Database
	CSRF     CSRF
	Session  Session
}

type Database struct {
	URL string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/zero?sslmode=disable"`
}

type CSRF struct {
	Token string `env:"CSRF_TOKEN"`
}

type Session struct {
	Key string `env:"SESSION_KEY"`
}
