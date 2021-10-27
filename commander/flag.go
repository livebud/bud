package commander

type Flag struct {
	Name  string
	Usage string
	value value
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

func (f *Flag) Bool(target *bool) *Bool {
	value := &Bool{target: target}
	f.value = &boolValue{inner: value}
	return value
}

func (f *Flag) verify(name string) error {
	return f.value.verify(name)
}

func verifyFlags(flags []*Flag) error {
	for _, flag := range flags {
		if err := flag.verify(flag.Name); err != nil {
			return err
		}
	}
	return nil
}
