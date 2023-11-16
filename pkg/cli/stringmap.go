package cli

import (
	"fmt"
	"strings"
)

type StringMap struct {
	target *map[string]string
	defval *map[string]string // default value
}

func (v *StringMap) Default(value map[string]string) {
	v.defval = &value
}

func (v *StringMap) Optional() {
	v.defval = new(map[string]string)
}

type stringMapValue struct {
	inner *StringMap
	set   bool
}

func (v *stringMapValue) verify(displayName string) error {
	if v.set {
		return nil
	} else if v.inner.defval != nil {
		*v.inner.target = *v.inner.defval
		return nil
	}
	return fmt.Errorf("missing %s", displayName)
}

func (v *stringMapValue) Set(val string) error {
	kv := strings.SplitN(val, ":", 2)
	if len(kv) != 2 {
		return fmt.Errorf("invalid key:value pair for %q", val)
	}
	if *v.inner.target == nil {
		*v.inner.target = map[string]string{}
	}
	(*v.inner.target)[kv[0]] = kv[1]
	v.set = true
	return nil
}

func (v *stringMapValue) String() string {
	if v.inner == nil {
		return ""
	} else if v.set {
		return v.format(*v.inner.target)
	} else if v.inner.defval != nil {
		return v.format(*v.inner.defval)
	}
	return ""
}

// Format as a string
func (v *stringMapValue) format(kv map[string]string) (out string) {
	i := 0
	for k, v := range kv {
		if i > 0 {
			out += " "
		}
		out += k + ":" + v
		i++
	}
	return out
}
