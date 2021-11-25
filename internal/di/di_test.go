package di_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/parser"

	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"gitlab.com/mnm/bud/internal/txtar"
	"gitlab.com/mnm/bud/vfs"
)

func goRun(dir string) (string, string, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", "run", "-mod=mod", "main.go")
	cmd.Env = os.Environ()
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	err := cmd.Run()
	if stderr.Len() > 0 {
		return "", stderr.String(), nil
	}
	if err != nil {
		return "", "", err
	}
	return stdout.String(), "", nil
}

// Searcher that bud uses
// - {importPath}
// - internal/{base(importPath)}
// - /{base(importPath)}
func searcher(modfile *mod.File) func(importPath string) (searchPaths []string) {
	return func(importPath string) (searchPaths []string) {
		modpath := modfile.ModulePath()
		base := path.Base(importPath)
		searchPaths = []string{
			importPath,
			path.Join(modpath, "internal", base),
			path.Join(modpath, base),
		}
		return searchPaths
	}
}

func runDI(dir, testscript string) (string, error) {
	modfile, err := mod.Default().Find(dir)
	if err != nil {
		return "", err
	}
	parser := parser.New(mod.Default())
	injector := di.New(modfile, parser)
	injector.Searcher = searcher(modfile)
	input, err := ioutil.ReadFile(filepath.Join(dir, "input.json"))
	if err != nil {
		return "", err
	}
	var in di.GenerateInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", fmt.Errorf("unable to unmarshal input.json: %w", err)
	}
	provider, err := injector.Generate(&in)
	if err != nil {
		return "", err
	}
	code := provider.File("Load")
	outDir := filepath.Join(dir, in.Target)
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return "", err
	}
	outPath := filepath.Join(outDir, "di.go")
	err = ioutil.WriteFile(outPath, []byte(code), 0644)
	if err != nil {
		return "", err
	}
	stdout, stderr, err := goRun(dir)
	if err != nil {
		return "", fmt.Errorf("error running generated code: %s", stderr)
	}
	if stderr != "" {
		return stderr, nil
	}
	return stdout, nil
}

func format(s string) string {
	// 8 spaces should be replaced with a tab
	s = strings.ReplaceAll(s, "        ", "\t")
	// Turn a tab into 2 spaces
	s = strings.ReplaceAll(s, "\t", "  ")
	// Remove all surrounding whitespace
	s = strings.TrimSpace(s)
	return s
}

func Test(t *testing.T) {
	is := is.New(t)
	des, err := ioutil.ReadDir("testdata")
	is.NoErr(err)
	for _, de := range des {
		de := de
		if de.IsDir() {
			continue
		}
		name := de.Name()
		if name[0] == '_' || filepath.Ext(name) != ".txt" {
			continue
		}
		testPath := filepath.Join("testdata", name)
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			dir := t.TempDir()
			fsys, err := txtar.ParseFile(testPath)
			is.NoErr(err)
			err = vfs.Write(dir, fsys)
			is.NoErr(err)
			expect, err := ioutil.ReadFile(filepath.Join(dir, "expect.txt"))
			is.NoErr(err)
			actual, err := runDI(dir, testPath)
			if err != nil {
				diff.TestString(t, format(string(expect)), err.Error())
				return
			}
			diff.TestString(t, format(string(expect)), format(actual))
		})
	}
}
