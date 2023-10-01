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

func (discard) Debug(args ...interface{}) error {
	return nil
}

func (discard) Debugf(msg string, args ...interface{}) error {
	return nil
}

func (discard) Info(args ...interface{}) error {
	return nil
}

func (discard) Infof(msg string, args ...interface{}) error {
	return nil
}

func (discard) Notice(args ...interface{}) error {
	return nil
}

func (discard) Noticef(msg string, args ...interface{}) error {
	return nil
}

func (discard) Warn(args ...interface{}) error {
	return nil
}

func (discard) Warnf(msg string, args ...interface{}) error {
	return nil
}

func (discard) Error(args ...interface{}) error {
	return nil
}

func (discard) Errorf(msg string, args ...interface{}) error {
	return nil
}
