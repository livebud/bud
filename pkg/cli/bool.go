package cli

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

type OptionalBool struct {
	target **bool
	defval *bool // default value
}

func (v *OptionalBool) Default(value bool) {
	v.defval = &value
}

type optionalBoolValue struct {
	inner *OptionalBool
	set   bool
}

var _ value = (*optionalBoolValue)(nil)

// IsBoolFlag allows --flag to be an alias for --flag true
func (v *optionalBoolValue) IsBoolFlag() bool {
	return true
}

func (v *optionalBoolValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = v.inner.defval
		return nil
	}
	return nil
}

func (v *optionalBoolValue) Set(val string) error {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*v.inner.target = &b
	v.set = true
	return nil
}

func (v *optionalBoolValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return strconv.FormatBool(**v.inner.target)
	} else if v.inner.defval != nil {
		return strconv.FormatBool(*v.inner.defval)
	}
	return ""
}
