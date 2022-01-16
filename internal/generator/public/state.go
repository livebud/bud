package public

type State struct {
	Embed bool
	Files []*File
}

type File struct {
	Path string
	Root string
}
