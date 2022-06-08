package golden

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/hexops/autogold"
)

func State(t *testing.T, state interface{}) {
	t.Helper()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("golden: unable to marshal state. %s", err)
	}
	autogold.Equal(t, autogold.Raw(string(data)), autogold.Dir(filepath.Join("testdata", "state")))
}

func Code(t *testing.T, code []byte) {
	t.Helper()
	autogold.Equal(t, autogold.Raw(code), autogold.Dir(filepath.Join("testdata", "code")))
}
