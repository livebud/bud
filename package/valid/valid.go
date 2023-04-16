package valid

import (
	"io/fs"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Dir validates that the name matches a valid directory
func Dir(name string) bool {
	return !invalidDir(name)
}

// Invalid dir check
func invalidDir(name string) bool {
	return len(name) == 0 || // Empty string
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		name == "bud" || // Named bud (reserved)
		strings.ToLower(name) != name // Has uppercase letters
}

// invalidFile check
func invalidFile(name string) bool {
	return len(name) == 0 || // Empty string
		name[0] == '_' || // Starts with _
		name[0] == '.' // Starts with .
}

// PluginDir validates that the name is a valid plugin directory
func PluginDir(name string) bool {
	return len(name) == 0 || // Empty string
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		strings.ToLower(name) != name // Has uppercase letters
}

// View validates that name matches a valid view
func View(name string) bool {
	return !invalidView(name)
}

// Invalid view check
func invalidView(name string) bool {
	ext := path.Ext(name)
	return len(name) == 0 || // Empty string
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		name == "bud" || // Named bud (reserved)
		ext == "" || // No extension
		ext == ".go" || // Go files
		name == "go.sum" || name == "go.mod" || // Go mod files
		name == "package.json" || // Node package file
		unicode.IsUpper(firstRune(name)) // Starts with a capital letter
}

// Invalid Go file
func invalidGoFile(name string) bool {
	return len(name) == 0 || // Empty string
		path.Ext(name) != ".go" ||
		name[0] == '_' || // Starts with _
		name[0] == '.' || // Starts with .
		name == "bud.go" || // Named bud (reserved)
		strings.HasSuffix(name, "_test") // Test file
}

// Invalid public file
func invalidPublicFile(name string) bool {
	return len(name) == 0 || // Empty string
		path.Ext(name) == "" ||
		name[0] == '_' || // Starts with _
		name[0] == '.' // Starts with .
}

func ControllerFile(name string) bool {
	return !invalidGoFile(name)
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

func CommandFile(name string) bool {
	return !invalidGoFile(name)
}

func GoFile(name string) bool {
	return !invalidGoFile(name)
}

func PublicFile(name string) bool {
	return !invalidPublicFile(name)
}

func WalkDirFunc(fn fs.WalkDirFunc) fs.WalkDirFunc {
	return func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check directories
		if de.IsDir() {
			if path != "." && invalidDir(de.Name()) {
				return fs.SkipDir
			}
			return fn(path, de, err)
		}
		// Check files
		if invalidFile(de.Name()) {
			return nil
		}
		return fn(path, de, err)
	}
}
