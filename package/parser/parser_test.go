package parser_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/package/conjure"
	"github.com/livebud/bud/package/merged"

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
	is.Equal(imp, "std/net/http")
	stct = def.Package().Struct("Request")
	is.True(stct != nil)
	is.Equal(stct.Name(), "Request")
}

func TestGenerate(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = `v0.0.8`
	is.NoErr(td.Write(ctx))
	cfs := conjure.New()
	merged := merged.Merge(os.DirFS(dir), cfs)
	cfs.GenerateFile("hello/hello.go", func(file *conjure.File) error {
		file.Data = []byte(`
			package hello
			import plugin "github.com/livebud/bud-test-plugin"
			type A struct { plugin.Answer }
		`)
		return nil
	})
	module, err := gomod.Find(dir)
	is.NoErr(err)
	is.Equal(module.Directory(), dir)
	p := parser.New(merged, module)
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
