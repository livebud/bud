package genfs

type EmbedFile struct {
	Data []byte
}

var _ FileGenerator = (*EmbedFile)(nil)

func (e *EmbedFile) GenerateFile(fsys FS, file *File) error {
	file.Data = e.Data
	return nil
}
