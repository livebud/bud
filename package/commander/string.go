package commander

import "fmt"

type String struct {
	target *string
	defval *string // default value
}

func (v *String) Default(value string) {
	v.defval = &value
}

func (v *String) Optional() {
	v.defval = new(string)
}

type stringValue struct {
	inner *String
	set   bool
}

func (v *stringValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = *v.inner.defval
		return nil
	}
	return fmt.Errorf("missing %s", displayName)
}

func (v *stringValue) Set(val string) error {
	*v.inner.target = val
	v.set = true
	return nil
}

func (v *stringValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return *v.inner.target
	} else if v.inner.defval != nil {
		return *v.inner.defval
	}
	return ""
}
