package log

// Filter logs by level
func Filter(level Level, handler Handler) Handler {
	return &filter{level, handler}
}

type filter struct {
	level   Level
	handler Handler
}

func (f *filter) Log(entry *Entry) error {
	if entry.Level < f.level {
		return nil
	}
	return f.handler.Log(entry)
}
