package viewkey

import (
	"fmt"
	"path"
	"strings"

	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/mux/ast"
)

func Infer(route string) (string, error) {
	if route == "" {
		return "", fmt.Errorf("empty route")
	} else if route == "/" {
		return "index", nil
	}
	route = trimExtension(route)
	r, err := mux.Parse(route)
	if err != nil {
		return "", err
	}
	lastSection := r.Sections[len(r.Sections)-1]
	action := "index"
	if _, ok := lastSection.(ast.Slot); ok {
		action = "show"
		r.Sections = r.Sections[:len(r.Sections)-1]
	} else if path, ok := lastSection.(*ast.Path); ok {
		if path.Value == "new" {
			action = "new"
			r.Sections = r.Sections[:len(r.Sections)-1]
		} else if path.Value == "edit" {
			action = "edit"
			r.Sections = r.Sections[:len(r.Sections)-1]
		} else if path.Value == "layout" {
			action = "layout"
			r.Sections = r.Sections[:len(r.Sections)-1]
		} else if path.Value == "frame" {
			action = "frame"
			r.Sections = r.Sections[:len(r.Sections)-1]
		} else if path.Value == "error" {
			action = "error"
			r.Sections = r.Sections[:len(r.Sections)-1]
		}
	}
	r = removeSlots(r)
	route = trimSlash(r.String())
	return path.Join(route, action), nil
}

func trimExtension(route string) string {
	return strings.TrimSuffix(route, path.Ext(route))
}

func removeSlots(r *ast.Route) *ast.Route {
	var remaining []ast.Section
	for _, s := range r.Sections {
		if _, ok := s.(ast.Slot); !ok {
			remaining = append(remaining, s)
		}
	}
	r.Sections = remaining
	return r
}

func trimSlash(route string) string {
	return strings.TrimLeft(route, "/")
}
