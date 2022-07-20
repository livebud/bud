// Package prompter is a small state machine for manipulating the terminal
// output during `bud run`.
//
// The public methods must be called in the right order:
//   1. Init
//   2. Reloading
//   3. SuccessReload or FailReload or NoReload
//   4. Reloading
//
// So the cursor can moves and clear lines properly. The behavior of the cursor
// across the different states looks like the following:
//
// | Ready         | Reloading            | Ready Again   |
// | :------------ | :------------------- | :------------ |
// | $ Ready on... | $ _Reloading...      | $ Ready on... |
// |   _           |   (move cursor up)   |   _           |
//
// TODO: find a better name for this package. Prompter sounds like it's reading
// user input from stdin, whereas this package is about managing `bud run` state
// and updating the terminal accordingly.
package prompter

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/watcher"
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
	events    []watcher.Event
	oldEvents []watcher.Event

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

func different(events, oldEvents []watcher.Event) bool {
	if len(events) != len(oldEvents) {
		return true
	}
	for i := range events {
		if events[i] != oldEvents[i] {
			return true
		}
	}
	return false
}

// Prompt successful reloads including time (in ms) and total times in a row.
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
	if different(p.events, p.oldEvents) {
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
func (p *Prompter) Reloading(events []watcher.Event) {
	if err := p.handleState(reload); err != nil {
		return
	}

	// Sort incoming events for faster comparison later between old and new events
	sort.Slice(events, func(i, j int) bool {
		return events[i].String() < events[j].String()
	})

	// Update events
	p.oldEvents = p.events
	p.events = events

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
func (p Prompter) NoReload() {
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
