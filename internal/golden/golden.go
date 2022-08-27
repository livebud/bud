package golden

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/matthewmueller/diff"

	"golang.org/x/tools/txtar"
)

var shouldUpdate = flag.Bool("update", false, "update golden files")

func TestGenerator(t testing.TB, state interface{}, code []byte) {
	t.Helper()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("golden: unable to marshal state. %s", err)
	}
	Test(t, &txtar.Archive{
		Files: []txtar.File{
			{
				Name: "state.json",
				Data: data,
			},
			{
				Name: "code.txt",
				Data: code,
			},
		},
	})
}

func Test(t testing.TB, actual *txtar.Archive) {
	t.Helper()
	filename := filepath.Join("testdata", t.Name()+".golden")
	formatted := txtar.Format(actual)
	expected, err := os.ReadFile(filename)
	if err != nil {
		if len(formatted) == 0 {
			return
		}
		expected = []byte("")
	}
	diff := difference(string(expected), string(formatted))
	if diff == "" {
		return
	}
	if *shouldUpdate {
		if err := writeFile(filename, formatted); err != nil {
			t.Fatalf("golden: unable to write golden file %s. %s", filename, err)
		}
	}
	t.Fatalf("golden: %s has unexpected changes.\n%s", filename, diff)
}

func writeFile(name string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	return os.WriteFile(name, data, 0644)
}

// TestString diffs two strings
func difference(expected string, actual string) string {
	if expected == actual {
		return ""
	}
	var b bytes.Buffer
	b.WriteString("\n\x1b[4mExpected\x1b[0m:\n")
	b.WriteString(expected)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mActual\x1b[0m: \n")
	b.WriteString(actual)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mDifference\x1b[0m: \n")
	b.WriteString(diff.String(expected, actual))
	b.WriteString("\n")
	return b.String()
}
