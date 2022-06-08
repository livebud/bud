/*
	Terminal manipulation based on ANSI Escape Sequences.
	Reference: https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
			 https://en.wikipedia.org/wiki/ANSI_escape_code

	Promptings must be placed in the right order
	Init -> Reloading -> Sucess Reload or Fail Reload -> Reloading -> ...

	So cursor moving and clear lines work properly ("_" is current cursor position):
		Reloading...               | -> _Reloading...   | -> Ready in 100ms (x23)
		_                          |                    |    _
		(current cursor position)  |  (move cursor up)  | (clear line and write new)
*/
package prompter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/livebud/bud/package/log/console"
)

// For prompting messages in the terminal
type Prompter struct {
	Counter int
	Error   string

	// For calculating total time
	startTime time.Time

	// Prevent overriding stdout and stderr
	StdOut    bytes.Buffer
	StdErr    bytes.Buffer
	oldStdOut bytes.Buffer
	oldStdErr bytes.Buffer

	// Use to prevent overriding error messages while compiling
	previousIsErr bool
}

// Clear line with cursor on.
func clearLine() {
	fmt.Print("\033[0K")
}

// Move cursor up 1 line.
func moveCursorUp() {
	fmt.Print("\033[1A")
}

func (p *Prompter) startTimer() {
	p.startTime = time.Now()
}

func (p *Prompter) Init() {
	p.Counter = 0
	fmt.Println("")
}

// Prompt failed reloads. Reset counter.
func (p *Prompter) FailReload(err string) {
	p.Counter = 0 // Reset counter
	p.previousIsErr = true
	console.Error(err)
}

// Prompt sucessful reloads including time (in ms) and total times in a row.
// Increase counter.
// Example: Ready in 100ms (x23).
func (p *Prompter) SuccessReload() {
	p.previousIsErr = false
	p.Counter++ // Increase counter

	// ! Temporary. Prevent overriding some unexpected errors we couldn't catch.
	if p.canOverridePreviousPrompt() {
		moveCursorUp()
		clearLine()
	}

	console.Info(fmt.Sprintf("Ready in %dms (x%d)", time.Since(p.startTime).Milliseconds(), p.Counter))
}

func (p *Prompter) canOverridePreviousPrompt() bool {
	newContentInStdOut := p.StdOut.String() != p.oldStdOut.String()
	newContentInStdErr := p.StdErr.String() != p.oldStdErr.String()

	// Renew
	p.oldStdOut = p.StdOut
	p.oldStdErr = p.StdErr

	// Only return true if there are nothing new on both stdout and stderr,
	// and no error previously
	return !newContentInStdErr && !newContentInStdOut && !p.previousIsErr
}

// Prompt "Reloading..." message.
// Start timer.
func (p *Prompter) Reloading() {
	p.startTimer() // For displaying reloading time in successful reloads

	// If there are anything new in stdout or stderr, we must not move cursor up and clear line
	// since it will probably override its messages.
	// Example if there's a error (we don't want to override):
	// Error: This is a error -> _Error: This is a error -> Reloading...
	// _                                                    _
	if p.canOverridePreviousPrompt() {
		moveCursorUp()
		clearLine()
	}
	console.Info("Reloading...")
}
