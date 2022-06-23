package command

// Flags is a shared struct that's passed between the bud CLI to the project CLI.
type Flag struct {
	Embed  bool
	Minify bool
	Hot    string
}
