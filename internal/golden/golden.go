package golden

import (
	"encoding/json"
	"testing"

	"github.com/hexops/autogold"
)

func State(t *testing.T, v interface{}) {
	t.Helper()
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("golden: unable to marshal state. %s", err)
	}
	autogold.Equal(t, autogold.Raw(string(buf)))
}
