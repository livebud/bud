package controller_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livebud/bud/framework/controller"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/parser"
)

func TestCheck(t *testing.T) {
	is := is.New(t)
	des, err := os.ReadDir("testdata")
	is.NoErr(err)
	for _, de := range des {
		name := strings.TrimSuffix(strings.TrimPrefix(de.Name(), "Test"), ".golden")
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			buf, err := os.ReadFile(filepath.Join("testdata", de.Name()))
			is.NoErr(err)
			state := new(controller.State)
			is.NoErr(json.Unmarshal(buf, state))
			code, err := controller.Generate(state)
			is.NoErr(err)
			if err := parser.Check(code); err != nil {
				fmt.Println(string(code))
				is.NoErr(err)
			}
		})
	}
}
