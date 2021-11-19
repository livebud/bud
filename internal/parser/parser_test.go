package parser_test

import (
	"path/filepath"
	"testing"

	"gitlab.com/mnm/bud/internal/modcache"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/internal/txtar"
	"gitlab.com/mnm/bud/vfs"
)

func TestFieldLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/field-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	p := parser.New(mod.New(modCache))
	pkg, err := p.Parse(filepath.Join(dir, "app", "hello"))
	is.NoErr(err)
	is.Equal(pkg.Name(), "hello")
	stct := pkg.Struct("A")
	is.True(stct != nil)
	field := stct.Field("S")
	is.True(field != nil)
	def, err := field.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Struct")
	pkg = def.Package()
	is.Equal(pkg.Name(), "module")
	stct = pkg.Struct("Struct")
	is.True(stct != nil)
	field = stct.Field("Dep")
	is.True(field != nil)
	def, err = field.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Dep")
	pkg = def.Package()
	is.Equal(pkg.Name(), "inner")
	stct = pkg.Struct("Dep")
	is.True(stct != nil)
}
