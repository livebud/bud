package gobin

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Run calls `go run -mod=mod main.go ...`
func Run(ctx context.Context, dir, mainpath string, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", "-mod=mod", mainpath}, args...)...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = dir
	// Setup a stderr pipe
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}
	// Process stderr output
	if err := processStderr(stderr); err != nil {
		return err
	}
	if err = cmd.Wait(); err != nil {
		if isCleanExit(err) {
			return nil
		}
		return err
	}
	return nil
}

// Interpret exit code 3 as an error-free exit. This allows interrupts like
// SIGINT to exit cleanly.
func isCleanExit(err error) bool {
	if e, ok := err.(*exec.ExitError); ok {
		if e.ExitCode() == 3 {
			return true
		}
	}
	return false
}

// Process stderr output
func processStderr(rc io.ReadCloser) error {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "go: found ") {
			continue
		}
		if strings.Contains(line, "exit status ") {
			continue
		}
		os.Stderr.WriteString(line + "\n")
	}
	return scanner.Err()
}
