package subtest

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Run a subtest as a child process.
func Run(t testing.TB, parent func(t testing.TB, cmd *exec.Cmd), child func(t testing.TB)) {
	if value := os.Getenv("CHILD"); value != "" {
		child(t)
		return
	}
	var args []string
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-test.count=") ||
			strings.HasPrefix(arg, "-test.v") ||
			strings.HasPrefix(arg, "-test.run") {
			continue
		}
		args = append(args, arg)
	}
	cmd := exec.Command(os.Args[0], append(args, "-test.v=true", "-test.run=^"+t.Name()+"$")...)
	cmd.Env = append(os.Environ(), "CHILD=1")
	parent(t, cmd)
}
