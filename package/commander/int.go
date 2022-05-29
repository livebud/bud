package commander

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

func (v *Int) Optional() {
	v.defval = new(int)
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
