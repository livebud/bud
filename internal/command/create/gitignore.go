package create

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed gitignore.gotext
var gitignore string

func (c *Command) generateGitIgnore(ctx context.Context, dir string) error {
	generator, err := gotemplate.Parse(".gitignore", gitignore)
	if err != nil {
		return err
	}
	var state struct{}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), code, 0644); err != nil {
		return err
	}
	return nil
}
