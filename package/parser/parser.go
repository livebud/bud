package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"unicode"

	"github.com/livebud/bud/package/gomod"
)

// New Go parser.
func New(fsys fs.FS, module *gomod.Module) *Parser {
	return &Parser{
		fsys:   fsys,
		module: module,
	}
}

// Parser for parsing Go code.
type Parser struct {
	fsys   fs.FS
	module *gomod.Module
}

// Parse a dir containing Go files.
func (p *Parser) Parse(dir string) (*Package, error) {
	imported, err := p.Import(dir)
	if err != nil {
		return nil, err
	}
	parsedPackage := &ast.Package{
		Name:  imported.Name,
		Files: make(map[string]*ast.File),
	}
	fset := token.NewFileSet()
	// Parse each valid Go file
	for _, filename := range imported.GoFiles {
		filename = path.Join(dir, filename)
		code, err := fs.ReadFile(p.fsys, filename)
		if err != nil {
			return nil, err
		}
		parsedFile, err := parser.ParseFile(fset, filename, code, parser.DeclarationErrors)
		if err != nil {
			return nil, err
		}
		parsedPackage.Files[filename] = parsedFile
	}
	pkg := newPackage(dir, p, p.module, parsedPackage)
	return pkg, nil
}

// Import the package, taking into account build tags and file name conventions
func (p *Parser) Import(dir string) (*build.Package, error) {
	return Import(p.fsys, dir)
}

func Import(fsys fs.FS, dir string) (*build.Package, error) {
	// TODO: figure out how to set the import path correctly to have better error
	// messages
	imported, err := buildContext(fsys).Import(".", dir, build.ImportMode(0))
	if err != nil {
		return nil, fmt.Errorf("parser: unable to import package %q. %w", dir, err)
	}
	return imported, nil
}

// A Context specifies the supporting context for a build. We mostly use the
// default context, but we want to override some of the values. This should be
// kept in sync with the keys in *build.Context
func buildContext(fsys fs.FS) *build.Context {
	context := build.Default
	return &build.Context{
		GOARCH:        context.GOARCH,
		GOOS:          context.GOOS,
		GOROOT:        context.GOROOT,
		GOPATH:        context.GOPATH,
		Dir:           context.Dir,
		CgoEnabled:    context.CgoEnabled,
		UseAllFiles:   context.UseAllFiles,
		Compiler:      context.Compiler,
		BuildTags:     context.BuildTags,
		ToolTags:      context.ToolTags,
		ReleaseTags:   context.ReleaseTags,
		InstallSuffix: context.InstallSuffix,

		// TODO: I'm not sure how to turn this into a call that uses the virtual
		// filesystem. It does rely on some os-specific filepath calls, but hasn't
		// seemed to affect the output.
		HasSubdir: context.HasSubdir,

		// // JoinPath joins the sequence of path fragments into a single path.
		// // If JoinPath is nil, Import uses filepath.Join.
		JoinPath: func(elem ...string) string {
			return path.Join(elem...)
		},

		// // SplitPathList splits the path list into a slice of individual paths.
		// // If SplitPathList is nil, Import uses filepath.SplitList.
		SplitPathList: filepath.SplitList,

		// // IsAbsPath reports whether path is an absolute path.
		// // If IsAbsPath is nil, Import uses filepath.IsAbs.
		IsAbsPath: func(elem string) bool {
			return path.IsAbs(elem)
		},

		// IsDir reports whether the path names a directory.
		// If IsDir is nil, Import calls os.Stat and uses the result's IsDir method.
		IsDir: func(path string) bool {
			fi, err := fs.Stat(fsys, path)
			if err != nil {
				// Error handling follows what build.Default does
				return false
			}
			return fi.IsDir()
		},

		// ReadDir returns a slice of fs.FileInfo, sorted by Name,
		// describing the content of the named directory.
		// If ReadDir is nil, Import uses ioutil.ReadDir.
		ReadDir: func(dir string) (fis []fs.FileInfo, err error) {
			des, err := fs.ReadDir(fsys, dir)
			if err != nil {
				return nil, err
			}
			for _, de := range des {
				fi, err := de.Info()
				if err != nil {
					return nil, err
				}
				fis = append(fis, fi)
			}
			return fis, nil
		},

		// OpenFile opens a file (not a directory) for reading.
		// If OpenFile is nil, Import uses os.Open.
		OpenFile: func(path string) (io.ReadCloser, error) {
			file, err := fsys.Open(path)
			if err != nil {
				return nil, err
			}
			return file, nil
		},
	}
}

func isPrivate(name string) bool {
	return unicode.IsLower(rune(name[0]))
}
