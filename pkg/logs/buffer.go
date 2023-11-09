package logs

import (
	"fmt"
	"sync"
)

func Buffer() *buffer {
	return &buffer{}
}

type buffer struct {
	mu      sync.Mutex
	entries []*Entry
}

var _ Handler = (*buffer)(nil)
var _ Log = (*buffer)(nil)

func (h *buffer) Log(entry *Entry) error {
	h.entries = append(h.entries, entry)
	return nil
}

// Entries returns a copy of the entries in the buffer
func (b *buffer) Entries() []*Entry {
	b.mu.Lock()
	defer b.mu.Unlock()
	entries := make([]*Entry, len(b.entries))
	copy(entries, b.entries)
	return entries
}

func (b *buffer) Field(key string, value interface{}) Log {
	return &sublogger{b, Fields{key: value}}
}

func (b *buffer) Fields(fields map[string]interface{}) Log {
	return &sublogger{b, fields}
}

func (b *buffer) log(lvl Level, args []interface{}, fields Fields) error {
	b.mu.Lock()
	b.entries = append(b.entries, createEntry(lvl, sprint(args...), fields))
	b.mu.Unlock()
	return nil
}

func (b *buffer) logf(lvl Level, msg string, args []interface{}, fields Fields) error {
	b.mu.Lock()
	b.entries = append(b.entries, createEntry(lvl, fmt.Sprintf(msg, args...), fields))
	b.mu.Unlock()
	return nil
}

func (b *buffer) Debug(args ...interface{}) error {
	return b.log(LevelDebug, args, nil)
}

func (b *buffer) Debugf(msg string, args ...interface{}) error {
	return b.logf(LevelDebug, msg, args, nil)
}

func (b *buffer) Info(args ...interface{}) error {
	return b.log(LevelInfo, args, nil)
}

func (b *buffer) Infof(msg string, args ...interface{}) error {
	return b.logf(LevelInfo, msg, args, nil)
}

func (b *buffer) Notice(args ...interface{}) error {
	return b.log(LevelNotice, args, nil)
}

func (b *buffer) Noticef(msg string, args ...interface{}) error {
	return b.logf(LevelNotice, msg, args, nil)
}

func (b *buffer) Warn(args ...interface{}) error {
	return b.log(LevelWarn, args, nil)
}

func (b *buffer) Warnf(msg string, args ...interface{}) error {
	return b.logf(LevelWarn, msg, args, nil)
}

func (b *buffer) Error(args ...interface{}) error {
	return b.log(LevelError, args, nil)
}

func (b *buffer) Errorf(msg string, args ...interface{}) error {
	return b.logf(LevelError, msg, args, nil)
}
