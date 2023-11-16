package cli

import "fmt"

type String struct {
	target *string
	defval *string // default value
}

func (v *String) Default(value string) {
	v.defval = &value
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

type OptionalString struct {
	target **string
	defval *string // default value
}

func (v *OptionalString) Default(value string) {
	v.defval = &value
}

type optionalStringValue struct {
	inner *OptionalString
	set   bool
}

var _ value = (*optionalStringValue)(nil)

func (v *optionalStringValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = v.inner.defval
		return nil
	}
	return nil
}

func (v *optionalStringValue) Set(val string) error {
	*v.inner.target = &val
	v.set = true
	return nil
}

func (v *optionalStringValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return **v.inner.target
	} else if v.inner.defval != nil {
		return *v.inner.defval
	}
	return ""
}
