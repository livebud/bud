package genfs

type Embed struct {
	Data []byte
}

var _ FileGenerator = (*Embed)(nil)

func (e *Embed) GenerateFile(fsys FS, file *File) error {
	file.Data = e.Data
	return nil
}
