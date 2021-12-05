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

	"gitlab.com/mnm/bud/internal/modcache"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/parser"

	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"gitlab.com/mnm/bud/internal/txtar"
	"gitlab.com/mnm/bud/vfs"
)

func goRun(cacheDir, appDir string) (string, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", "run", "-mod=mod", "main.go")
	cmd.Env = append(os.Environ(), "GOMODCACHE="+cacheDir, "GOPRIVATE=*")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = appDir
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
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
	module := mod.Default()
	modFile, err := module.Find(dir)
	if err != nil {
		return "", err
	}
	parser := parser.New(module)
	injector := di.New(parser)
	injector.Searcher = searcher(modFile)
	input, err := ioutil.ReadFile(filepath.Join(dir, "input.json"))
	if err != nil {
		return "", err
	}
	var in di.GenerateInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", fmt.Errorf("unable to unmarshal input.json: %w", err)
	}
	// Add the project modfile to each initial dependency
	for _, dep := range in.Dependencies {
		if dep.ModFile == nil {
			dep.ModFile = modFile
		}
	}
	provider, err := injector.Generate(&in)
	if err != nil {
		return "", err
	}
	code := provider.File("Load")
	fmt.Println(code)
	// TODO: provide a modFile method for doing this.
	// Right now ResolveDirectory also stats the final dir
	targetDir := modFile.Directory(strings.TrimPrefix(in.Target, modFile.ModulePath()))
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return "", err
	}
	outPath := filepath.Join(targetDir, "di.go")
	err = ioutil.WriteFile(outPath, []byte(code), 0644)
	if err != nil {
		return "", err
	}
	stdout, err := goRun(modcache.Default().Directory(), dir)
	if err != nil {
		return "", fmt.Errorf("error running generated code: %s", err)
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

// TODO: This should be merged with the rest of the tests
func runDIWithModules(modDir, appDir, testscript string) (string, error) {
	modCache := modcache.New(modDir)
	err := modCache.Write(modcache.Modules{
		"mod.test/three@v1.0.0": modcache.Files{
			"inner/inner.go": `
				package inner

				import (
					"fmt"
				)

				type Three struct {}

				func (t Three) String() string {
					return fmt.Sprintf("Three{}")
				}
			`,
		},
		"mod.test/two@v0.0.1": modcache.Files{
			"struct.go": `
				package two

				type Struct struct {
				}
			`,
		},
		"mod.test/two@v0.0.2": modcache.Files{
			"go.mod": `
				module mod.test/two

				require (
					mod.test/three v1.0.0
				)
			`,
			"struct.go": `
				package two

				import (
					"mod.test/three/inner"
					"fmt"
				)

				type Two struct {
					inner.Three
				}

				func (t *Two) String() string {
					return fmt.Sprintf("&Two{Three: %s}", t.Three)
				}
			`,
		},
	})
	if err != nil {
		return "", err
	}
	module := mod.New(modCache)
	modFile, err := module.Find(appDir)
	if err != nil {
		return "", err
	}
	parser := parser.New(module)
	injector := di.New(parser)
	injector.Searcher = searcher(modFile)
	input, err := ioutil.ReadFile(filepath.Join(appDir, "input.json"))
	if err != nil {
		return "", err
	}
	var in di.GenerateInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", fmt.Errorf("unable to unmarshal input.json: %w", err)
	}
	// Add the project modfile to each initial dependency
	for _, dep := range in.Dependencies {
		if dep.ModFile == nil {
			dep.ModFile = modFile
		}
	}
	provider, err := injector.Generate(&in)
	if err != nil {
		return "", err
	}
	code := provider.File("Load")
	// TODO: provide a modFile method for doing this.
	// Right now ResolveDirectory also stats the final dir
	targetDir := modFile.Directory(strings.TrimPrefix(in.Target, modFile.ModulePath()))
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return "", err
	}
	outPath := filepath.Join(targetDir, "di.go")
	err = ioutil.WriteFile(outPath, []byte(code), 0644)
	if err != nil {
		return "", err
	}
	stdout, err := goRun(modDir, appDir)
	if err != nil {
		return "", fmt.Errorf("error running generated code: %s", err)
	}
	return stdout, nil
}

func TestModuleNested(t *testing.T) {
	is := is.New(t)
	des, err := ioutil.ReadDir("testdata")
	is.NoErr(err)
	for _, de := range des {
		de := de
		if de.IsDir() {
			continue
		}
		name := de.Name()
		if name != "_module_nested.txt" {
			continue
		}
		testPath := filepath.Join("testdata", name)
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			modDir := t.TempDir()
			appDir := t.TempDir()
			fsys, err := txtar.ParseFile(testPath)
			is.NoErr(err)
			err = vfs.Write(appDir, fsys)
			is.NoErr(err)
			expect, err := ioutil.ReadFile(filepath.Join(appDir, "expect.txt"))
			is.NoErr(err)
			actual, err := runDIWithModules(modDir, appDir, testPath)
			if err != nil {
				diff.TestString(t, format(string(expect)), err.Error())
				return
			}
			diff.TestString(t, format(string(expect)), format(actual))
		})
	}
}
