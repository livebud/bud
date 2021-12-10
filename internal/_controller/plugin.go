package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	// Embed templates
	_ "embed"

	"github.com/go-duo/duoc/compiler"
	"github.com/go-duo/duoc/di"
	"github.com/go-duo/duoc/internal/go/imports"
	"github.com/go-duo/duoc/go/parser"
	"github.com/go-duo/duoc/template"
	"gitlab.com/mnm/duo/log"
)

// New fn
func New(injector *di.Injector, log log.Log, parser *parser.Parser) *Plugin {
	return &Plugin{injector, log, parser}
}

// Plugin struct
type Plugin struct {
	injector *di.Injector
	log      log.Log
	parser   *parser.Parser
}

var _ compiler.Plugin = (*Plugin)(nil)

// ID of the dst
func (p *Plugin) ID() string {
	return "controller"
}

var viewPath = filepath.Join("generated", "view", "view.go")

// Find files
func (p *Plugin) Find(fs compiler.FS, graph compiler.Graph) error {
	files, err := fs.Glob("controller/{*,**}.go")
	if err != nil {
		return err
	}
	viewExists, err := fs.Exists(viewPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		genpath := filepath.Join("generated", filepath.Dir(file), "controller.go")
		graph.Link(genpath, file)
		if viewExists {
			graph.Link(genpath, viewPath)
		}
	}
	return nil
}

// Delete the web server
func (p *Plugin) Delete(node compiler.Node) error {
	p.log.Info("deleting %q", node.Path())
	return os.RemoveAll(filepath.Dir(node.Path()))
}

// Create the web server
func (p *Plugin) Create(node compiler.Node) error {
	p.log.Info("creating %q", node.Path())
	return p.Upsert(node)
}

// Update the web server
func (p *Plugin) Update(node compiler.Node) error {
	p.log.Info("updating %q", node.Path())
	return p.Upsert(node)
}

//go:embed template.gotext
var generate string
var generator = template.MustParse("controller", generate)

// Upsert the controller
func (p *Plugin) Upsert(target compiler.Node) error {
	state, err := p.compile(target)
	if err != nil {
		return err
	}
	return generator.Overwrite(target.Path(), state)
}

// Load the controller state
func (p *Plugin) compile(target compiler.Node) (state *Controller, err error) {
	source, err := p.findController(target)
	if err != nil {
		return nil, err
	}
	pkg, err := p.parser.Parse(filepath.Dir(source.Path()))
	if err != nil {
		return nil, err
	}
	modfile, err := pkg.Modfile()
	if err != nil {
		return nil, err
	}
	// Load the state for the current run
	loader := &loader{
		injector: p.injector,
		imports:  imports.New(),
		contexts: newContextMap(),
		modfile:  modfile,
		target:   target,
		view:     p.findView(target),
	}
	return loader.load(pkg)
}

func newContextMap() *contextSet {
	return &contextSet{map[string]*Context{}}
}

type contextSet struct {
	contextMap map[string]*Context
}

func (c *contextSet) Add(context *Context) {
	c.contextMap[context.Function] = context
}

func (c *contextSet) List() (contexts []*Context) {
	for _, context := range c.contextMap {
		contexts = append(contexts, context)
	}
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Function < contexts[j].Function
	})
	return contexts
}

// Find the source controller
func (p *Plugin) findController(target compiler.Node) (controller compiler.Node, err error) {
	dir := compiler.CommonDir(target)
	targetDir := filepath.Join(dir, "generated")
	targetPath, err := filepath.Rel(targetDir, target.Path())
	if err != nil {
		return nil, err
	}
	prereqs := target.Prerequisites()
	for _, prereq := range prereqs {
		sourcePath, err := filepath.Rel(dir, prereq.Path())
		if err != nil {
			return nil, err
		}
		if sourcePath == targetPath {
			return prereq, nil
		}
	}
	return nil, fmt.Errorf("controller: unable to find source node for %q", target.Path())
}

// Find the source controller
func (p *Plugin) findView(target compiler.Node) (view compiler.Node) {
	viewPath := filepath.Join("generated", "view", "view.go")
	for _, prereq := range target.Prerequisites() {
		if strings.HasSuffix(prereq.Path(), viewPath) {
			return prereq
		}
	}
	return nil
}
