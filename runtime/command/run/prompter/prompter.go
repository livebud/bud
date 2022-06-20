/*
	Terminal manipulation based on ANSI Escape Sequences.
	Reference: https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
			 https://en.wikipedia.org/wiki/ANSI_escape_code

	Promptings must be placed in the right order
		Init -> Reloading -> `Sucess reload` or `Fail reload` or `Made no reload` ->
 			Reloading -> ...

	So cursor moving and clear lines work properly ("_" is current cursor position):
		Reloading...               | -> _Reloading...   | -> Ready in 100ms (x23)
		_                          |                    |    _
		(current cursor position)  |  (move cursor up)  | (clear line and write new)
*/
package prompter

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/livebud/bud/package/log/console"
)

// States
const (
	fail     = "fail"
	success  = "success"
	reload   = "reload"
	noReload = "no re"
)

// For prompting messages in the terminal
type Prompter struct {
	Counter int

	// For calculating total time
	startTime time.Time

	// Prevent overriding stdout and stderr
	StdOut    bytes.Buffer
	StdErr    bytes.Buffer
	oldStdOut bytes.Buffer
	oldStdErr bytes.Buffer

	// Path to changed files
	paths    []string
	oldPaths []string

	// For storing states (fail, success, reload,...)
	state    string
	oldState string

	reloadMessage    string
	listeningAddress string
}

// Clear line with cursor on.
func clearLine() {
	fmt.Fprint(os.Stderr, "\033[0K")
}

// Move cursor up 1 line.
func moveCursorUp() {
	fmt.Fprint(os.Stderr, "\033[1A")
}

func (p *Prompter) startTimer() {
	p.startTime = time.Now()
}

// Init -> Reloading -> `Sucess reload` or `Fail reload` or `Made no reload` ->
// 		Reloading -> ...
var nextState = map[string]string{
	"init":   "reload",
	reload:   "fail/success/no re",
	fail:     "reload",
	success:  "reload",
	noReload: "reload",
}

// Ensure all states are in proper arrangement
func (p *Prompter) handleState(state string) error {
	// Update state
	p.oldState = p.state
	p.state = state

	switch {
	case p.state == p.oldState:
		return fmt.Errorf("duplicated state: %s", p.state)
	case !strings.Contains(nextState[p.oldState], p.state):
		return fmt.Errorf("invalid state, expected %s instead of %s", nextState[p.oldState], p.state)
	}

	return nil
}

func (p *Prompter) Init(listeningAddress string) {
	p.state = "init"
	p.Counter = 0 // Init counter
	p.listeningAddress = listeningAddress
}

func (p *Prompter) blankStdOut() (result bool) {
	result = p.StdOut.String() == p.oldStdOut.String()

	// Update stdout
	p.oldStdOut = p.StdOut

	return result
}

func (p *Prompter) blankStdErr() (result bool) {
	result = p.StdErr.String() == p.oldStdErr.String()

	// Update stderr
	p.oldStdErr = p.StdErr

	return result
}

// Prompt failed reloads. Reset counter.
func (p *Prompter) FailReload(err string) {
	if err := p.handleState(fail); err != nil {
		return
	}
	p.Counter = 0 // Reset counter
	console.Error(err)
}

func different(paths, oldPaths []string) bool {
	if len(paths) != len(oldPaths) {
		return true
	}

	sort.Strings(paths)
	sort.Strings(oldPaths)
	for i := range paths {
		if paths[i] != oldPaths[i] {
			return true
		}
	}

	return false
}

// Prompt sucessful reloads including time (in ms) and total times in a row.
// Increase counter.
// Example: Ready on http://127.0.0.1:3000 in 264ms (x141)
func (p *Prompter) SuccessReload() {
	if err := p.handleState(success); err != nil {
		return
	}

	p.Counter++ // Increase counter

	// Prevent override
	if p.blankStdErr() && p.blankStdOut() {
		moveCursorUp()
		clearLine()
	}

	// Reset counter if user changed working file
	if different(p.paths, p.oldPaths) {
		p.Counter = 1
	}

	p.reloadMessage = fmt.Sprintf("Ready on %s in %dms (x%d)",
		p.listeningAddress,
		time.Since(p.startTime).Milliseconds(),
		p.Counter,
	)

	console.Info(p.reloadMessage)
}

// Prompt "Reloading..." message.
// Start timer.
func (p *Prompter) Reloading(paths []string) {
	if err := p.handleState(reload); err != nil {
		return
	}

	// Update paths
	p.oldPaths = p.paths
	p.paths = paths

	// Prevent override
	if p.blankStdErr() && p.blankStdOut() && p.oldState != fail {
		moveCursorUp()
		clearLine()
	}

	// For displaying reloading time in successful reloads
	p.startTimer()

	console.Info("Reloading...")
}

// Don't change reload message currently in the terminal
func (p Prompter) MadeNoReload() {
	// Doesn't check for error, so it could be used multiple times in a row.
	p.handleState(noReload)

	// Prevent override
	if p.blankStdErr() && p.blankStdOut() {
		moveCursorUp()
		clearLine()
	}

	if p.reloadMessage != "" {
		console.Info(p.reloadMessage)
	} else {
		console.Info("Listening on " + p.listeningAddress)
	}
}