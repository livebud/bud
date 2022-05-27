package commander

import (
	"fmt"
	"strconv"
)

type Bool struct {
	target *bool
	defval *bool // default value
}

func (v *Bool) Default(value bool) {
	v.defval = &value
}

func (v *Bool) Optional() {
	v.defval = new(bool)
}

type boolValue struct {
	inner *Bool
	set   bool
}

func (v *boolValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = *v.inner.defval
		return nil
	}
	return fmt.Errorf("missing %s", displayName)
}

func (v *boolValue) Set(val string) (err error) {
	*v.inner.target, err = strconv.ParseBool(val)
	if err != nil {
		return err
	}
	v.set = true
	return nil
}

func (v *boolValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return strconv.FormatBool(*v.inner.target)
	} else if v.inner.defval != nil {
		return strconv.FormatBool(*v.inner.defval)
	}
	return "false"
}

// IsBoolFlag allows --flag to be an alias for --flag true
func (v *boolValue) IsBoolFlag() bool {
	return true
}
