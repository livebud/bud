package cli

import (
	"fmt"
	"strings"
)

type Strings struct {
	target   *[]string
	defval   *[]string // default value
	optional bool
}

func (v *Strings) Default(values ...string) {
	v.defval = &values
}

type stringsValue struct {
	inner *Strings
	set   bool
}

func (v *stringsValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = *v.inner.defval
		return nil
	} else if v.inner.optional {
		return nil
	}
	return fmt.Errorf("missing %s", displayName)
}

func (v *stringsValue) Set(val string) error {
	*v.inner.target = append(*v.inner.target, val)
	v.set = true
	return nil
}

func (v *stringsValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return strings.Join(*v.inner.target, ", ")
	} else if v.inner.defval != nil {
		return strings.Join(*v.inner.defval, ", ")
	}
	return ""
}
