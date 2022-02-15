package overlay

import "gitlab.com/mnm/bud/package/conjure"

type FileGenerator interface {
	GenerateFile(fsys F, file *File) error
}

type File struct {
	*conjure.File
}
