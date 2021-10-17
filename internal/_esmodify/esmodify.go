package esmodify

import (
	"bytes"
	"fmt"
	"regexp"
)

var reImport = regexp.MustCompile(`([A-Z_a-z$][A-Z_a-z0-9]*)?\(?"(__LIVEBUD_EXTERNAL__:([^"]+))"\)?`)
var importBytes = []byte(`import`)

// This function rewrites require statements and updates the path on imports
func ReplaceDependencyPaths(content []byte) []byte {
	identifiers := map[string]bool{}
	out := new(bytes.Buffer)
	code := new(bytes.Buffer)
	since := 0
	// Submatches: [
	//  (0) matchStart,
	//  (1) matchEnd,
	//  (2) requireOrImportStart,
	//  (3) requireOrImportEnd,
	//  (4) modulePathStart,
	//  (5) modulePathEnd,
	//  (6) moduleNameStart,
	//  (7) moduleNameEnd,
	// ]
	for _, submatches := range reImport.FindAllSubmatchIndex(content, -1) {
		// Write the bytes since the last match
		code.Write(content[since:submatches[0]])
		// Update since with the end of the match
		since = submatches[1]
		// Get the path of the node module
		path := string(content[submatches[6]:submatches[7]])
		// Handle require(...) or import(...)
		var importOrRequire []byte
		if submatches[2] >= 0 && submatches[3] >= 0 {
			importOrRequire = content[submatches[2]:submatches[3]]
		}
		// We have a require(...), replace the whole expression
		if importOrRequire != nil && !bytes.Equal(importOrRequire, importBytes) {
			identifier := "__" + toIdentifier(path) + "$"
			code.WriteString(identifier)
			// Only add this import if we haven't seen this identifier yet
			if !identifiers[identifier] {
				out.WriteString(importStatement(identifier, path))
				identifiers[identifier] = true
			}
			continue
		}
		// Otherwise, we'll just replace the path
		code.Write(content[submatches[0]:submatches[4]])
		code.WriteString("/bud/node_modules/" + path)
		code.Write(content[submatches[5]:submatches[1]])
	}
	// Write the remaining bytes
	code.Write(content[since:])
	// Write code to out
	out.Write(code.Bytes())
	return out.Bytes()
}

func toIdentifier(importPath string) string {
	p := []byte(importPath)
	for i, c := range p {
		switch c {
		case '/', '-', '@', '.':
			p[i] = '_'
		default:
			p[i] = c
		}
	}
	return string(p)
}

func importStatement(identifier, name string) string {
	return fmt.Sprintf(`import %s from "/bud/node_modules/%s"`+"\n", identifier, name)
}
