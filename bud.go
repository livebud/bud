package bud

type Env string

const (
	Development Env = "development"
	Test        Env = "test"
	Preview     Env = "preview"
	Production  Env = "production"
)
