/*
	Terminal manipulation based on ANSI Escape Sequences.
	Reference: https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
			 https://en.wikipedia.org/wiki/ANSI_escape_code

	Promptings must be placed in the right order
		PromptReloading -> PromptSucessReload/PromptFailedReload -> PromptReloading -> ...

	So cursor moving and clear lines work probably ("_" is current cursor position):
		Reloading...				-> _Reloading        -> Ready in 100ms (x23)
		_										    _
		(current cursor position)	 (move cursor up)    (clear line and write new)
*/
package terminal

import (
	"fmt"
	"time"

	"github.com/livebud/bud/package/log/console"
)

// Clear line with cursor on.
func clearLine() {
	fmt.Print("\033[0K")
}

// Move cursor up 1 line.
func moveCursorUp() {
	fmt.Print("\033[1A")
}

// Prompt failed reloads. Reset counter. Can only be used after PromptReloading.
func PromptFailedReload(err string, counter *int) {
	*counter = 0 // Reset counter
	console.Error(err)
	console.Error("") // Prevent overriding error messages
}

// Prompt sucessfully reloads including time (in ms) and total sucessful reloads in a row.
// Example: Ready in 100ms (x23). Can only be used after PromptReloading
func PromptSucessReload(totalTime time.Duration, counter *int) {
	*counter++ // Increase counter
	moveCursorUp()
	clearLine()
	console.Info(fmt.Sprintf("Ready in %dms (x%d)", totalTime.Milliseconds(), *counter))
}

// Prompt "Reloading..." message. Need to know whether user wrote anything in
// the stdout to prevent overrides. Must be used before PromptSucessReload or
// PromptFailedReload
func PromptReloading(userPrinted bool) {
	// If user printed anything out, we must not move cursor up and clear line
	// since it will probably override their print messages
	if !userPrinted {
		moveCursorUp()
		clearLine()
	}
	console.Info("Reloading...")
}
