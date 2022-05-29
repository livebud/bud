package scaffold

import (
	"go/format"
	"path"
	"path/filepath"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/vfs"
	"golang.org/x/sync/errgroup"
)

type Templates []*Template

func (templates Templates) Write(fsys vfs.ReadWritable) error {
	eg := new(errgroup.Group)
	for _, template := range templates {
		template := template
		eg.Go(func() error { return template.Write(fsys) })
	}
	return eg.Wait()

}

type Template struct {
	Path  string
	Code  string
	State interface{}
}

func (t *Template) Write(fsys vfs.ReadWritable) error {
	generator, err := gotemplate.Parse(t.Path, t.Code)
	if err != nil {
		return err
	}
	code, err := generator.Generate(t.State)
	if err != nil {
		return err
	}
	// Format Go code automatically
	if filepath.Ext(t.Path) == ".go" {
		code, err = format.Source(code)
		if err != nil {
			return err
		}
	}
	if err := fsys.MkdirAll(path.Dir(t.Path), 0755); err != nil {
		return err
	}
	if err := fsys.WriteFile(t.Path, code, 0644); err != nil {
		return err
	}
	return nil
}
