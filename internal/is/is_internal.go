package is

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var nocolor = false

func init() {
	// prefer https://no-color.org (with any value)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		nocolor = true
	}
}

const (
	colorNormal  = "\u001b[39m"
	colorComment = "\u001b[31m"
	colorFile    = "\u001b[90m"
	colorType    = "\u001b[90m"
)

type T interface {
	FailNow()
}

func New(t T) *I {
	return &I{os.Stdout, !nocolor, t, t.FailNow, map[string]struct{}{}}
}

type I struct {
	writer   io.Writer
	colorful bool
	t        T
	fail     func()
	helpers  map[string]struct{} // functions to be skipped when writing file/line info
}

func (is *I) log(args ...interface{}) {
	s := is.decorate(fmt.Sprint(args...))
	fmt.Fprintf(is.writer, s)
	is.fail()
}

func (is *I) logf(format string, args ...interface{}) {
	is.log(fmt.Sprintf(format, args...))
}

func (is *I) valWithType(v interface{}) string {
	if isNil(v) {
		return "<nil>"
	}
	if is.colorful {
		return fmt.Sprintf("%[1]s%[3]T(%[2]s%[3]v%[1]s)%[2]s", colorType, colorNormal, v)
	}
	return fmt.Sprintf("%[1]T(%[1]v)", v)
}

var reArgs = regexp.MustCompile(`\$\{[0-9]+\}`)

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
// this function was copied from the testing framework and modified.
func (is *I) decorate(s string) string {
	path, lineNumber, ok := is.callerinfo() // decorate + log + public function.
	file := filepath.Base(path)
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		lineNumber = 1
	}
	buf := new(bytes.Buffer)
	if is.colorful {
		buf.WriteString(colorFile)
	}
	fmt.Fprintf(buf, "%s:%d: ", file, lineNumber)
	if is.colorful {
		buf.WriteString(colorNormal)
	}

	s = escapeFormatString(s)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	args, hasArgs := loadArguments(path, lineNumber)
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		// expand arguments (if /\$[0-9]+/ is present)
		line = reArgs.ReplaceAllStringFunc(line, func(old string) string {
			n, err := strconv.Atoi(strings.Trim(old, "${}"))
			if err != nil {
				return old
			}
			if hasArgs && len(args) >= n {
				return args[n-1]
			}
			return old
		})
		buf.WriteString(line)
	}
	buf.WriteString("\n")
	return buf.String()
}

// callerName gives the function name (qualified with a package path)
// for the caller after skip frames (where 0 means the current function).
func callerName(skip int) string {
	// Make room for the skip PC.
	var pc [1]uintptr
	n := runtime.Callers(skip+2, pc[:]) // skip + runtime.Callers + callerName
	if n == 0 {
		panic("is: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

// The maximum number of stack frames to go through when skipping helper functions for
// the purpose of decorating log messages.
const maxStackLen = 50

var reIsSourceFile = regexp.MustCompile(`is(_internal)?\.go$`)

func (is *I) callerinfo() (path string, line int, ok bool) {
	var pc [maxStackLen]uintptr
	// Skip two extra frames to account for this function
	// and runtime.Callers itself.
	n := runtime.Callers(2, pc[:])
	if n == 0 {
		panic("is: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	var firstFrame, frame runtime.Frame
	for more := true; more; {
		frame, more = frames.Next()
		if reIsSourceFile.MatchString(frame.File) {
			continue
		}
		if firstFrame.PC == 0 {
			firstFrame = frame
		}
		if _, ok := is.helpers[frame.Function]; ok {
			// Frame is inside a helper function.
			continue
		}
		return frame.File, frame.Line, true
	}
	// If no "non-helper" frame is found, the first non is frame is returned.
	return firstFrame.File, firstFrame.Line, true
}

// areEqual gets whether a equals b or not.
func areEqual(a, b interface{}) bool {
	if isNil(a) && isNil(b) {
		return true
	}
	if isNil(a) || isNil(b) {
		return false
	}
	if reflect.DeepEqual(a, b) {
		return true
	}
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	return aValue == bValue
}

// isNil gets whether the object is nil or not.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

// escapeFormatString escapes strings for use in formatted functions like Sprintf.
func escapeFormatString(fmt string) string {
	return strings.Replace(fmt, "%", "%%", -1)
}

// loadArguments gets the arguments from the function call
// on the specified line of the file.
func loadArguments(path string, line int) ([]string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	i := 1
	for s.Scan() {
		if i != line {
			i++
			continue
		}
		text := s.Text()
		braceI := strings.Index(text, "(")
		if braceI == -1 {
			return nil, false
		}
		text = text[braceI+1:]
		cs := bufio.NewScanner(strings.NewReader(text))
		cs.Split(bufio.ScanBytes)
		j := 0
		c := 1
		for cs.Scan() {
			switch cs.Text() {
			case ")":
				c--
			case "(":
				c++
			}
			if c == 0 {
				break
			}
			j++
		}
		text = text[:j]
		return strings.Split(text, ", "), true
	}
	return nil, false
}

func (is *I) formatArgs(args []interface{}) string {
	switch len(args) {
	case 0:
		return ""
	case 1:
		b := new(strings.Builder)
		b.WriteString(". ")
		if is.colorful {
			b.WriteString(colorComment)
		}
		b.WriteString(fmt.Sprint(args[0]))
		if is.colorful {
			b.WriteString(colorNormal)
		}
		return b.String()
	default:
		b := new(strings.Builder)
		b.WriteString(". ")
		if is.colorful {
			b.WriteString(colorComment)
		}
		b.WriteString(fmt.Sprintf(fmt.Sprint(args[0]), args[1:]...))
		if is.colorful {
			b.WriteString(colorNormal)
		}
		return b.String()
	}
}

// containsElement try loop over the list check if the list includes the element.
// return (false, false) if impossible.
// return (true, false) if element was not found.
// return (true, true) if element was found.
//
// Pulled from: https://github.com/stretchr/testify/blob/083ff1c0449867d0d8d456483ee5fab8e0c0e1e6/assert/assertions.go#L721
func containsElement(list interface{}, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	listType := reflect.TypeOf(list)
	if listType == nil {
		return false, false
	}
	listKind := listType.Kind()
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()
	if listKind == reflect.String {
		elementValue := reflect.ValueOf(element)
		return true, strings.Contains(listValue.String(), elementValue.String())
	}
	if listKind == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}
	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false
}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}
	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
