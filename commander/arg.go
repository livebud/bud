package commander

type Arg struct {
	Name  string
	Usage string
	value value
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

func (a *Arg) Strings(target *[]string) *Strings {
	value := &Strings{target: target}
	a.value = &stringsValue{inner: value}
	return value
}

func (a *Arg) verify(name string) error {
	return a.value.verify(name)
}

func verifyArgs(args []*Arg) error {
	for _, arg := range args {
		if err := arg.verify(arg.Name); err != nil {
			return err
		}
	}
	return nil
}
