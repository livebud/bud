package testdir

import (
	"os"
	"path/filepath"

	"github.com/lithammer/dedent"
	"golang.org/x/sync/errgroup"
)

func WriteFiles(dir string, files map[string]string) error {
	eg := new(errgroup.Group)
	if _, ok := files["go.mod"]; !ok {
		files["go.mod"] = `module example.com
		go 1.16
		`
	}
	for name, contents := range files {
		name, contents := name, contents
		path := filepath.Join(dir, name)
		eg.Go(func() error {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(path, []byte(dedent.Dedent(contents)), 0644); err != nil {
				return err
			}
			return nil
		})
	}
	return eg.Wait()
}
