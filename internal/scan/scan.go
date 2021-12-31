package scan

type Scanner interface {
	Scan() bool
	Err() error
	Text() string
}
