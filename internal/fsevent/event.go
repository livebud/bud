package fsevent

import "strings"

// Event describes a set of file operations.
type Event uint32

// These are the generalized file operations that can trigger a notification.
const (
	create Event = 1 << iota
	update
	delete
	rename
)

func (event Event) String() string {
	// Use a buffer for efficient string concatenation
	buffer := new(strings.Builder)
	if event&create == create {
		buffer.WriteString("|create")
	}
	if event&delete == delete {
		buffer.WriteString("|delete")
	}
	if event&update == update {
		buffer.WriteString("|update")
	}
	if event&rename == rename {
		buffer.WriteString("|rename")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:] // Strip leading pipe
}

// Create event
func Create(path string) (Event, []byte) {
	return create, []byte(path)
}

// Update event
func Update(path string) (Event, []byte) {
	return update, []byte(path)
}

// Delete event
func Delete(path string) (Event, []byte) {
	return delete, []byte(path)
}

// Rename event
func Rename(path string) (Event, []byte) {
	return rename, []byte(path)
}
