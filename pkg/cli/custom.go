package cli

import (
	"fmt"
)

type Custom struct {
	target func(s string) error
	defval *string // default value
}

func (v *Custom) Default(value string) {
	v.defval = &value
}

func (v *Custom) Optional() {
	v.defval = new(string)
}

type customValue struct {
	inner *Custom
	set   bool
}

var _ value = (*customValue)(nil)

func (v *customValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		return v.inner.target(*v.inner.defval)
	}
	return fmt.Errorf("missing %s", displayName)
}

func (v *customValue) Set(val string) error {
	err := v.inner.target(val)
	v.set = true
	return err
}

func (v *customValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return ""
	} else if v.inner.defval != nil {
		return *v.inner.defval
	}
	return ""
}
