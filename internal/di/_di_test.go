package di_test

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/parser"
)

func Test(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile, err := mod.Find(wd)
	is.NoErr(err)
	projectDir := modfile.Directory()
	testDir, err := filepath.Abs("_tmp")
	is.NoErr(err)
	err = os.RemoveAll(testDir)
	is.NoErr(err)
	err = os.MkdirAll(testDir, 0755)
	is.NoErr(err)
	p := testscript.Params{
		WorkdirRoot: testDir,
		Dir:         "testdata",
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"di": diCmd,
		},
		Setup: func(e *testscript.Env) error {
			e.Setenv("GONOSUMDB", "*")
			e.Setenv("DUOC_DIR", projectDir)
			return nil
		},
	}
	err = gotooltest.Setup(&p)
	is.NoErr(err)
	testscript.Run(t, p)
	t.Cleanup(func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(testDir))
		}
	})
}

type dependencies []string

func (d *dependencies) String() string {
	return "dependencies"
}

func (d *dependencies) Set(value string) error {
	*d = append(*d, value)
	return nil
}

type externals []string

func (e *externals) String() string {
	return "external data types"
}

func (e *externals) Set(value string) error {
	*e = append(*e, value)
	return nil
}

type Searcher struct {
	Modfile *mod.File
}

// Searcher that duo uses
// - {importPath}
// - internal/{base(importPath)}
// - /{base(importPath)}
func (s *Searcher) Search(importPath string) (searchPaths []string) {
	modpath := s.modfile.ModulePath()
	base := path.Base(importPath)
	searchPaths = []string{
		importPath,
		path.Join(modpath, "internal", base),
		path.Join(modpath, base),
	}
	return searchPaths
}

// di command
func diCmd(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("di: doesn't support !")
	}
	command := flag.NewFlagSet("di", flag.ExitOnError)
	var dependencies dependencies
	command.Var(&dependencies, "d", "dependencies")
	var externals externals
	command.Var(&externals, "e", "external data types")
	output := command.String("o", "", "output")
	hoist := command.Bool("hoist", false, "hoist the node")
	if err := command.Parse(args); err != nil {
		ts.Fatalf("di: error parsing flags %+s", err)
	} else if len(dependencies) == 0 {
		ts.Fatalf("di: missing dependency '-d=<import>.<type>'")
	} else if *output == "" {
		ts.Fatalf("di: missing output flag '-o=<path>'")
	}
	modfinder := mod.New()
	modfile, err := modfinder.Find(ts.MkAbs(""))
	if err != nil {
		ts.Fatalf("di: unable to find go.mod %+s", err)
	}
	parser := parser.New(mod.New())
	injector := di.New(modfile, parser)
	searcher := &Searcher{modfile}
	injector.Searcher = searcher.Search
	var input di.GenerateInput
	input.Target = path.Join(modfile.ModulePath(), filepath.Dir(*output))
	for _, dependency := range dependencies {
		input.Dependencies = append(input.Dependencies, toDependency(ts, modfile, dependency))
	}
	for _, external := range externals {
		input.Externals = append(input.Externals, toDependency(ts, modfile, external))
	}
	if *hoist {
		input.Hoist = true
	}
	// graph, err := injector.Print(&di.PrintInput{
	// 	Dependencies: input.Dependencies,
	// 	Externals:    input.Externals,
	// 	Hoist:        input.Hoist,
	// })
	// if err != nil {
	// 	ts.Fatalf("di: failed to print the graph")
	// }
	// fmt.Println(graph)
	// Use the searcher that duo will use
	provider, err := injector.Generate(&input)
	if err != nil {
		ts.Fatalf("di: wiring failed: %+s", err)
	}
	code := provider.File("Load")
	outdir := filepath.Dir(ts.MkAbs(*output))
	if err = os.MkdirAll(outdir, 0755); err != nil {
		ts.Fatalf("di: failed to make directory %s because %+s", outdir, err)
	}
	if err = ioutil.WriteFile(ts.MkAbs(*output), []byte(code), 0644); err != nil {
		ts.Fatalf("di: failed to write %s file %+s", *output, err)
	}
}

func toDependency(ts *testscript.TestScript, modfile *mod.File, dependency string) *di.Dependency {
	parts := strings.SplitN(dependency, ".", 2)
	if len(parts) < 2 {
		ts.Fatalf("di: external must have form '<import>.<type>'. got %q ", dependency)
	}
	// This should handle both stdlib (e.g. "net/http") and directories (e.g. "web")
	importPath := parts[0]
	if _, err := os.Stat(ts.MkAbs(parts[0])); err == nil {
		importPath, err = modfile.ResolveImport(ts.MkAbs(parts[0]))
		if err != nil {
			ts.Fatalf("di: unable to resolve import %s because %+s", parts[0], err)
		}
	}
	// Create the dependency
	dep := &di.Dependency{
		Import: importPath,
		Type:   parts[1],
	}
	return dep
}
