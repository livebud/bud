package cli

type Flag struct {
	name  string
	help  string
	short byte
	value value
}

func (f *Flag) Short(short byte) *Flag {
	f.short = short
	return f
}

func (f *Flag) Optional() *OptionalFlag {
	return &OptionalFlag{f}
}

func (f *Flag) Int(target *int) *Int {
	value := &Int{target: target}
	f.value = &intValue{inner: value}
	return value
}

func (f *Flag) String(target *string) *String {
	value := &String{target: target}
	f.value = &stringValue{inner: value}
	return value
}

func (f *Flag) Strings(target *[]string) *Strings {
	*target = []string{}
	value := &Strings{target: target}
	f.value = &stringsValue{inner: value}
	return value
}

func (f *Flag) StringMap(target *map[string]string) *StringMap {
	*target = map[string]string{}
	value := &StringMap{target: target}
	f.value = &stringMapValue{inner: value}
	return value
}

func (f *Flag) Bool(target *bool) *Bool {
	value := &Bool{target: target}
	f.value = &boolValue{inner: value}
	return value
}

// Custom allows you to define a custom parsing function
func (f *Flag) Custom(fn func(string) error) *Custom {
	value := &Custom{target: fn}
	f.value = &customValue{inner: value}
	return value
}

func (f *Flag) verify(name string) error {
	return f.value.verify("--" + name)
}

type OptionalFlag struct {
	f *Flag
}

func (f *OptionalFlag) String(target **string) *OptionalString {
	value := &OptionalString{target: target}
	f.f.value = &optionalStringValue{inner: value}
	return value
}

func (f *OptionalFlag) Int(target **int) *OptionalInt {
	value := &OptionalInt{target: target}
	f.f.value = &optionalIntValue{inner: value}
	return value
}

func (f *OptionalFlag) Bool(target **bool) *OptionalBool {
	value := &OptionalBool{target: target}
	f.f.value = &optionalBoolValue{inner: value}
	return value
}

func (f *OptionalFlag) Strings(target *[]string) *Strings {
	*target = []string{}
	value := &Strings{target: target, optional: true}
	f.f.value = &stringsValue{inner: value}
	return value
}

func verifyFlags(flags []*Flag) error {
	for _, flag := range flags {
		if err := flag.verify(flag.name); err != nil {
			return err
		}
	}
	return nil
}
