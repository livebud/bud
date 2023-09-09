package color

import "fmt"

func Ignore() Writer {
	return &ignore{}
}

type ignore struct{}

var _ Writer = (*ignore)(nil)

func (i *ignore) Enabled() bool {
	return false
}

func (i *ignore) Blue(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func (i *ignore) Red(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func (i *ignore) Yellow(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func (i *ignore) Dim(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func (i *ignore) Green(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func (i *ignore) Pink(v ...interface{}) string {
	return fmt.Sprint(v...)
}
