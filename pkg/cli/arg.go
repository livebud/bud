package cli

type Arg struct {
	name  string
	value value
}

func (a *Arg) Optional() *OptionalArg {
	return &OptionalArg{a}
}

func (a *Arg) Int(target *int) *Int {
	value := &Int{target: target}
	a.value = &intValue{inner: value}
	return value
}

func (a *Arg) String(target *string) *String {
	value := &String{target: target}
	a.value = &stringValue{inner: value}
	return value
}

func (a *Arg) StringMap(target *map[string]string) *StringMap {
	*target = map[string]string{}
	value := &StringMap{target: target}
	a.value = &stringMapValue{inner: value}
	return value
}

// Custom allows you to define a custom parsing function
func (a *Arg) Custom(fn func(string) error) *Custom {
	value := &Custom{target: fn}
	a.value = &customValue{inner: value}
	return value
}

func (a *Arg) verify(name string) error {
	return a.value.verify(name)
}

type OptionalArg struct {
	a *Arg
}

func (a *OptionalArg) String(target **string) *OptionalString {
	value := &OptionalString{target: target}
	a.a.value = &optionalStringValue{inner: value}
	return value
}

func (a *OptionalArg) Int(target **int) *OptionalInt {
	value := &OptionalInt{target: target}
	a.a.value = &optionalIntValue{inner: value}
	return value
}

func (a *OptionalArg) Bool(target **bool) *OptionalBool {
	value := &OptionalBool{target: target}
	a.a.value = &optionalBoolValue{inner: value}
	return value
}

func verifyArgs(args []*Arg) error {
	for _, arg := range args {
		if err := arg.verify("<" + arg.name + ">"); err != nil {
			return err
		}
	}
	return nil
}
