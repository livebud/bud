package parser_test

import (
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/internal/modcache"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/internal/txtar"
	"gitlab.com/mnm/bud/vfs"
)

func TestStructLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/struct-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	modFinder := mod.New(mod.WithCache(modCache))
	module, err := modFinder.Find(filepath.Join(dir, "app"))
	is.NoErr(err)
	p := parser.New(module)
	pkg, err := p.Parse("hello")
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
	is.Equal(pkg.Name(), "two")
	modFile := pkg.Module()
	is.Equal(modFile.Import(), "mod.test/two")
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
	modFile = pkg.Module()
	is.Equal(modFile.Import(), "mod.test/three")
}

func TestInterfaceLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/interface-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	modFinder := mod.New(mod.WithCache(modCache))
	module, err := modFinder.Find(filepath.Join(dir, "app"))
	is.NoErr(err)
	p := parser.New(module)
	pkg, err := p.Parse("hello")
	is.NoErr(err)
	is.Equal(pkg.Name(), "hello")
	stct := pkg.Struct("A")
	is.True(stct != nil)
	field := stct.Field("S")
	is.True(field != nil)
	def, err := field.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Interface")
	pkg = def.Package()
	is.Equal(pkg.Name(), "two")
	module = pkg.Module()
	is.Equal(module.Import(), "mod.test/two")
	iface := pkg.Interface("Interface")
	is.True(iface != nil)
	is.Equal(iface.Name(), "Interface")
	method := iface.Method("Test")
	is.True(method != nil)
	results := method.Results()
	is.Equal(len(results), 1)
	def, err = results[0].Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Interface")
	pkg = def.Package()
	is.Equal(pkg.Name(), "inner")
	iface = pkg.Interface("Interface")
	is.True(iface != nil)
	method = iface.Method("String")
	is.True(method != nil)
	is.Equal(method.Name(), "String")
	module = pkg.Module()
	is.Equal(module.Import(), "mod.test/three")
}

func TestAliasLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/alias-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	modFinder := mod.New(mod.WithCache(modCache))
	module, err := modFinder.Find(filepath.Join(dir, "app"))
	is.NoErr(err)
	p := parser.New(module)
	pkg, err := p.Parse(".")
	is.NoErr(err)
	is.Equal(pkg.Name(), "main")
	alias := pkg.Alias("Middleware")
	is.True(alias != nil)
	is.Equal(alias.Name(), "Middleware")
	def, err := alias.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Middleware")
	pkg = def.Package()
	is.Equal(pkg.Name(), "public")
	alias = pkg.Alias("Middleware")
	is.True(alias != nil)
	def, err = alias.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Interface")
	pkg = def.Package()
	is.Equal(pkg.Name(), "public")
	middleware := pkg.Interface("Interface")
	is.True(middleware != nil)
	method := middleware.Method("Middleware")
	is.True(method != nil)
	is.Equal(method.Name(), "Middleware")
}

type Module struct {
	Modules map[string]map[string]string
	Files   map[string]string
}

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
}

func makeModule(t *testing.T, m Module) *mod.Module {
	is := is.New(t)
	modCache := modcache.Default()
	if m.Modules != nil {
		cacheDir := t.TempDir()
		modCache = modcache.New(cacheDir)
		err := modCache.Write(m.Modules)
		is.NoErr(err)
	}
	appDir := t.TempDir()
	if m.Files != nil {
		for path, file := range m.Files {
			m.Files[path] = redent(file)
		}
		err := vfs.Write(appDir, vfs.Map(m.Files))
		is.NoErr(err)
	}
	modFinder := mod.New(mod.WithCache(modCache))
	module, err := modFinder.Find(appDir)
	is.NoErr(err)
	return module
}

func TestNetHTTP(t *testing.T) {
	module := makeModule(t, Module{
		Files: map[string]string{
			"go.mod": `module app.com/app`,
			"app.go": `
				package app

				import "net/http"

				type A struct {
					*http.Request
				}
			`,
		},
	})
	is := is.New(t)
	p := parser.New(module)
	pkg, err := p.Parse(".")
	is.NoErr(err)
	stct := pkg.Struct("A")
	is.True(stct != nil)
	is.Equal(stct.Name(), "A")
	field := stct.Field("Request")
	is.True(field != nil)
	is.Equal(field.Name(), "Request")
	def, err := field.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Request")
	pkg = def.Package()
	imp, err := pkg.Import()
	is.NoErr(err)
	is.Equal(imp, "std/net/http")
	stct = def.Package().Struct("Request")
	is.True(stct != nil)
	is.Equal(stct.Name(), "Request")
}
