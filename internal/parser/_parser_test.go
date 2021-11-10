package parser_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/parser"

	"github.com/matryer/is"
	"github.com/rogpeppe/go-internal/goproxytest"
	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

var _ = (*parser.Parser)(nil)

func Test(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile, err := mod.Load(wd)
	is.NoErr(err)
	budVersion := modfile.Version("gitlab.com/mnm/bud")
	is.True(len(budVersion) > 0)
	workDir, err := ioutil.TempDir("", "bud")
	is.NoErr(err)
	replacementDuo := modfile.Replacement("gitlab.com/mnm/bud")
	srv, err := goproxytest.NewServer(filepath.Join(wd, "testdata", "mod"), "")
	is.NoErr(err)
	// Create a new modcache each test run
	modCache, err := ioutil.TempDir("", "gomodcache")
	if err != nil {
		t.Fatalf("cannot get tempdir: %v", err)
	}
	p := testscript.Params{
		WorkdirRoot: workDir,
		Dir:         filepath.Join(wd, "testdata"),
		Cmds:        map[string]func(ts *testscript.TestScript, neg bool, args []string){},
		Condition: func(cond string) (bool, error) {
			switch cond {
			case "replace":
				return replacementDuo != "", nil
			case "ci":
				return os.Getenv("CI") != "", nil
			default:
				return false, fmt.Errorf("Unknown condition %q", cond)
			}
		},
		Setup: func(e *testscript.Env) error {
			e.Setenv("GONOSUMDB", "*")
			e.Setenv("GOMODCACHE", modCache)
			e.Setenv("GOPROXY", srv.URL+",https://proxy.golang.org,direct")
			e.Setenv("DUOC_DIR", modfile.Directory())
			e.Setenv("DUO_DIR", replacementDuo)
			e.Setenv("DUO_VERSION", budVersion)
			return nil
		},
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}
