package markdoc

import "github.com/livebud/bud/framework/transform2/transformrt"

func New(compiler *Markdoc) *Transform {
	return &Transform{compiler}
}

type Transform struct {
	md *Markdoc
}

type Markdoc struct {
}

func (c *Markdoc) Compile(md []byte) (html []byte, err error) {
	return []byte("<h1>" + string(md) + "</h1>"), nil
}

func (t *Transform) MdToSvelte(file *transformrt.File) error {
	html, err := t.md.Compile(file.Data)
	if err != nil {
		return err
	}
	file.Data = html
	return nil
}
