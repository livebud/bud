package goroot

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Find returns the GOROOT value, using either an explicitly
// provided environment variable, a GOROOT that contains the current
// os.Executable value, or else the GOROOT that the binary was built
// with from runtime.GOROOT().
//
// This function is a heavily modified version of the following:
// https://github.com/golang/go/blob/89044b6d423a07bea3b6f80210f780e859dd2700/src/cmd/go/internal/cfg/cfg.go#L369
func Find() string {
	if env := os.Getenv("GOROOT"); env != "" {
		return filepath.Clean(env)
	}
	def := ""
	if r := runtime.GOROOT(); r != "" {
		def = filepath.Clean(r)
	}
	if runtime.Compiler == "gccgo" {
		// gccgo has no real GOROOT, and it certainly doesn't
		// depend on the executable's location.
		return def
	}
	// Run `go env GOROOT` to compute the GOROOT.
	// TODO: this takes about 8ms on boot, see if we can find a better way to
	// compute the correct GOROOT.
	// Note: The tricky part is making it work with --trimpath.
	cmd := exec.Command("go", "env", "GOROOT")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return def
	}
	return strings.TrimSpace(stdout.String())
}
