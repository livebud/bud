package parser_test

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/parser"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/txtar"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/vfs"
)

// TODO: replace txtar with testdir
func TestStructLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/struct-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	appDir := filepath.Join(dir, "app")
	module, err := gomod.Find(appDir, gomod.WithModCache(modCache))
	is.NoErr(err)
	is.Equal(module.Import(), "app.com")
	p := parser.New(module, module)
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

// TODO: replace txtar with testdir
func TestInterfaceLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/interface-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	appDir := filepath.Join(dir, "app")
	module, err := gomod.Find(appDir, gomod.WithModCache(modCache))
	is.NoErr(err)
	p := parser.New(module, module)
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

// TODO: replace txtar with testdir
func TestAliasLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/alias-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	appDir := filepath.Join(dir, "app")
	module, err := gomod.Find(appDir, gomod.WithModCache(modCache))
	is.NoErr(err)
	p := parser.New(module, module)
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

func TestNetHTTP(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["app.go"] = `
		package app

		import "net/http"

		type A struct {
			*http.Request
		}
	`
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	p := parser.New(module, module)
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
	is.Equal(imp, "net/http")
	stct = def.Package().Struct("Request")
	is.True(stct != nil)
	is.Equal(stct.Name(), "Request")
}

func TestGenerate(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = `v0.0.8`
	is.NoErr(td.Write(ctx))
	fsys := genfs.New(dag.Discard, os.DirFS(dir), log)
	fsys.GenerateFile("hello/hello.go", func(fsys genfs.FS, file *genfs.File) error {
		file.Data = []byte(`
			package hello
			import plugin "github.com/livebud/bud-test-plugin"
			type A struct { plugin.Answer }
		`)
		return nil
	})
	module, err := gomod.Find(dir)
	is.NoErr(err)
	p := parser.New(fsys, module)
	// Parse a virtual package
	pkg, err := p.Parse("hello")
	is.NoErr(err)
	is.Equal(pkg.Name(), "hello")
	stct := pkg.Struct("A")
	is.True(stct != nil)
	is.Equal(stct.Name(), "A")
	field := stct.Field("Answer")
	is.True(field != nil)
	is.Equal(field.Name(), "Answer")
	// Visit real dependencies from the virtual package
	def, err := field.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Answer")
	pkg = def.Package()
	is.Equal(pkg.Name(), "plugin")
	importPath, err := pkg.Import()
	is.NoErr(err)
	is.Equal(importPath, "github.com/livebud/bud-test-plugin")
	alias := pkg.Alias("Answer")
	is.True(alias != nil)
	is.Equal(alias.Name(), "Answer")
	def, err = alias.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Answer")
	pkg = def.Package()
	is.Equal(pkg.Name(), "plugin")
	importPath, err = pkg.Import()
	is.NoErr(err)
	is.Equal(importPath, "github.com/livebud/bud-test-nested-plugin")
	alias = pkg.Alias("Answer")
	is.True(alias != nil)
	is.Equal(alias.Name(), "Answer")
}

func TestAliasLookupModule(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	budModule, err := gomod.Find(".")
	is.NoErr(err)
	dep := budModule.File().Require("github.com/livebud/transpiler")
	td.Modules["github.com/livebud/transpiler"] = dep.Version
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	fsys := virtual.Tree{
		"bud/transpiler/transpiler.go": &virtual.File{Data: []byte(`
			package transpiler
			import "app.com/runtime/transpiler"
			func New(tr *transpiler.Transpiler) {}
		`)},
		"runtime/transpiler/transpiler.go": &virtual.File{Data: []byte(`
			package transpiler
			import "github.com/livebud/transpiler"
			type Transpiler = transpiler.Transpiler
		`)},
	}
	p := parser.New(fsys, module)
	pkg, err := p.Parse("bud/transpiler")
	is.NoErr(err)
	is.Equal(pkg.Name(), "transpiler")
	newFn := pkg.Function("New")
	is.True(newFn != nil)
	params := newFn.Params()
	is.Equal(len(params), 1)
	is.Equal(params[0].Name(), "tr")
	def, err := params[0].Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Transpiler")
	is.Equal(def.Kind(), parser.KindStruct)
	pkg = def.Package()
	is.Equal(pkg.Name(), "transpiler")
	importPath, err := pkg.Import()
	is.NoErr(err)
	is.Equal(importPath, "app.com/runtime/transpiler")
	alias := def.Package().Alias("Transpiler")
	is.True(alias != nil)
	is.Equal(alias.Name(), "Transpiler")
	def, err = alias.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Transpiler")
	is.Equal(def.Kind(), parser.KindStruct)
	pkg = def.Package()
	is.Equal(pkg.Name(), "transpiler")
	importPath, err = pkg.Import()
	is.NoErr(err)
	is.Equal(importPath, "github.com/livebud/transpiler")
	is.Equal(pkg.Directory(), path.Join(module.ModCache(), "github.com/livebud/transpiler@"+dep.Version))
}
