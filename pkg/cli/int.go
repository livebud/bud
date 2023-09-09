package cli

import (
	"fmt"
	"strconv"
)

type Int struct {
	target *int
	defval *int
}

func (v *Int) Default(value int) {
	v.defval = &value
}

type intValue struct {
	inner *Int
	set   bool
}

func (v *intValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = *v.inner.defval
		return nil
	}
	return fmt.Errorf("missing %s", displayName)
}

func (v *intValue) Set(val string) error {
	n, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*v.inner.target = n
	v.set = true
	return nil
}

func (v *intValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return strconv.Itoa(*v.inner.target)
	} else if v.inner.defval != nil {
		return strconv.Itoa(*v.inner.defval)
	}
	return ""
}

type OptionalInt struct {
	target **int
	defval *int
}

func (v *OptionalInt) Default(value int) {
	v.defval = &value
}

type optionalIntValue struct {
	inner *OptionalInt
	set   bool
}

func (v *optionalIntValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = v.inner.defval
		return nil
	}
	return nil
}

func (v *optionalIntValue) Set(val string) error {
	n, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*v.inner.target = &n
	v.set = true
	return nil
}

func (v *optionalIntValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return strconv.Itoa(**v.inner.target)
	} else if v.inner.defval != nil {
		return strconv.Itoa(*v.inner.defval)
	}
	return ""
}
