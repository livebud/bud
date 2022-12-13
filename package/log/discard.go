package log

var Discard = discard{}

type discard struct{}

var _ Log = discard{}

func (d discard) Field(key string, value interface{}) Log {
	return d
}

func (d discard) Fields(fields map[string]interface{}) Log {
	return d
}

func (discard) Debug(msg string, args ...interface{}) error {
	return nil
}

func (discard) Info(msg string, args ...interface{}) error {
	return nil
}

func (discard) Notice(msg string, args ...interface{}) error {
	return nil
}

func (discard) Warn(msg string, args ...interface{}) error {
	return nil
}

func (discard) Error(msg string, args ...interface{}) error {
	return nil
}

func (discard) Err(err error, msg string, args ...interface{}) error {
	return nil
}
