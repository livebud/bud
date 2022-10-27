package stacktrace

import (
	"runtime"
	"strconv"
	"strings"
)

// Frames returns a stack of callers
func Frames(skip int) []Frame {
	var pcs [32]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	var st stack = pcs[0:n]
	return st.Frames()
}

// Source prints the most recent frame of the stack trace.
func Source(skip int) string {
	// Add 1 since we're now within the Source function
	frames := Frames(skip + 1)
	top := frames[0]
	str := new(strings.Builder)
	if name := top.Name(); name != "" {
		str.WriteString(name)
	} else if file := top.File(); file != "" {
		str.WriteString(file)
	}
	if line := top.Line(); line > 0 {
		str.WriteString(":" + strconv.Itoa(line))
	}
	return str.String()
}

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Frames() []Frame {
	f := make([]Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((*s)[i])
	}
	return f
}

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) File() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) Line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
func (f Frame) Name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}
	return fn.Name()
}
