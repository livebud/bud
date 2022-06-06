/*
	Terminal manipulation based on ANSI Escape Sequences.
	Reference: https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
			 https://en.wikipedia.org/wiki/ANSI_escape_code

	Promptings must be placed in the right order
	Init -> Start timer -> Reloading -> SucessReload/FailReload -> Reloading -> Start timer -> ...

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

	// Prevent overriding error messages while compiling
	previousIsErr bool
}

// Clear line with cursor on.
func clearLine() {
	// fmt.Print("\033[0K")
}

// Move cursor up 1 line.
func moveCursorUp() {
	// fmt.Print("\033[1A")
}

func (p *Prompter) Init() {
	p.Counter = 0
}

func (p *Prompter) StartTimer() {
	p.startTime = time.Now()
}

// Prompt failed reloads. Reset counter.
func (p *Prompter) FailReload(err string) {
	p.Counter = 0 // Reset counter
	p.previousIsErr = true
	console.Error(err)
}

// Prompt sucessfully reloads including time (in ms) and total sucessful reloads in a row.
// Increase counter.
// Example: Ready in 100ms (x23).
func (p *Prompter) SuccessReload() {
	p.Counter++ // Increase counter
	moveCursorUp()
	clearLine()
	console.Info(fmt.Sprintf("Ready in %dms (x%d)", time.Since(p.startTime).Milliseconds(), p.Counter))
}

func (p *Prompter) okToClearLine() bool {
	newContentInStdOut := p.StdOut.String() != p.oldStdOut.String()
	newContentInStdErr := p.StdErr.String() != p.oldStdErr.String()

	// Renew
	p.oldStdOut = p.StdOut
	p.oldStdErr = p.StdErr

	// Only return true if there are nothing new on both stdout and stderr
	return !newContentInStdErr && !newContentInStdOut
}

// Prompt "Reloading..." message.
func (p *Prompter) Reloading() {
	// If there are anything new in stdout or stderr, we must not move cursor up and clear line
	// since it will probably override its messages
	// Example if there's a error (we don't want to override):
	// Error: This is a error -> _Error: This is a error -> Reloading...
	// _                                                    _
	if p.okToClearLine() {
		moveCursorUp()
		clearLine()
	}
	console.Info("Reloading...")
}
