package genfs

import "strings"

// Event describes a set of file operations.
type Event uint32

// These are the generalized file operations that can trigger a notification.
const (
	CreateEvent Event = 1 << iota
	WriteEvent
	RemoveEvent
)

func (event Event) String() string {
	// Use a buffer for efficient string concatenation
	buffer := new(strings.Builder)

	if event&CreateEvent == CreateEvent {
		buffer.WriteString("|Create")
	}
	if event&RemoveEvent == RemoveEvent {
		buffer.WriteString("|Remove")
	}
	if event&WriteEvent == WriteEvent {
		buffer.WriteString("|Write")
	}
	// if event&RenameEvent == RenameEvent {
	// 	buffer.WriteString("|RenameEvent")
	// }
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:] // Strip leading pipe
}
