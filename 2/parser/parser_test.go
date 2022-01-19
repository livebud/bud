package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/mnm/bud/2/parser"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/modcache"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/2/mod"
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
	appDir := filepath.Join(dir, "app")
	module, err := mod.Find(appDir, mod.WithCache(modCache))
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

func TestInterfaceLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/interface-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	appDir := filepath.Join(dir, "app")
	module, err := mod.Find(appDir, mod.WithCache(modCache))
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

func TestAliasLookup(t *testing.T) {
	is := is.New(t)
	testfile, err := txtar.ParseFile("testdata/alias-lookup.txt")
	is.NoErr(err)
	dir := t.TempDir()
	err = vfs.Write(dir, testfile)
	is.NoErr(err)
	modCache := modcache.New(filepath.Join(dir, "mod"))
	appDir := filepath.Join(dir, "app")
	module, err := mod.Find(appDir, mod.WithCache(modCache))
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
	appDir := t.TempDir()
	err := vfs.Write(appDir, vfs.Map{
		"go.mod": []byte(`module app.com/app`),
		"app.go": []byte(`
			package app

			import "net/http"

			type A struct {
				*http.Request
			}
		`),
	})
	is.NoErr(err)
	module, err := mod.Find(appDir)
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
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(modcache.Modules{
		"mod.test/two@v0.0.1": modcache.Files{
			"go.mod":   "module mod.test/two",
			"const.go": "package two\ntype Answer = int",
		},
		"mod.test/two@v0.0.2": modcache.Files{
			"go.mod":   "module mod.test/two",
			"const.go": "package two\ntype Answer = string",
		},
		"mod.test/module@v1.2.3": modcache.Files{
			"go.mod":   "module mod.test/module",
			"const.go": "package module\ntype Answer = string",
		},
		"mod.test/module@v1.2.4": modcache.Files{
			"go.mod":   "module mod.test/module\nrequire mod.test/two v0.0.2",
			"const.go": "package module\nimport \"mod.test/two\"\ntype Answer = two.Answer",
		},
	})
	is.NoErr(err)
	appDir := t.TempDir()
	err = vfs.Write(appDir, vfs.Map{
		"go.mod": []byte("module app.com\nrequire mod.test/module v1.2.4"),
		"app.go": []byte("package app\nimport \"mod.test/module\"\nvar a = module.Answer"),
	})
	is.NoErr(err)
	genfs := gen.New(os.DirFS(appDir))
	genfs.Add(map[string]gen.Generator{
		"hello/hello.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("package hello\nimport \"mod.test/module\"\ntype A struct { module.Answer }"))
			return nil
		}),
	})
	module, err := mod.Find(appDir, mod.WithCache(modCache))
	is.NoErr(err)
	is.Equal(module.Directory(), appDir)
	p := parser.New(genfs, module)
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
	is.Equal(pkg.Name(), "module")
	importPath, err := pkg.Import()
	is.NoErr(err)
	is.Equal(importPath, "mod.test/module")
	alias := pkg.Alias("Answer")
	is.True(alias != nil)
	is.Equal(alias.Name(), "Answer")
	def, err = alias.Definition()
	is.NoErr(err)
	is.Equal(def.Name(), "Answer")
	pkg = def.Package()
	is.Equal(pkg.Name(), "two")
	importPath, err = pkg.Import()
	is.NoErr(err)
	is.Equal(importPath, "mod.test/two")
	alias = pkg.Alias("Answer")
	is.True(alias != nil)
	is.Equal(alias.Name(), "Answer")
}
