package commander

type Arg struct {
	Name  string
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

func verifyArgs(args []*Arg) error {
	for _, arg := range args {
		if err := arg.verify(arg.Name); err != nil {
			return err
		}
	}
	return nil
}

// func (a *Arg) Optional() *optionalArg {

// }

// type optionalArg struct {
// 	Name  string
// 	value value
// }

// func (a *optionalArg) Int(target **int) *Int {
// 	value := &Int{target: target}
// 	a.value = &intValue{inner: value}
// 	return value
// }

// func (a *optionalArg) String(target **string) *String {
// 	value := &String{target: target}
// 	a.value = &stringValue{inner: value}
// 	return value
// }

// func (a *optionalArg) Strings(target **[]string) *Strings {
// 	value := &Strings{target: target}
// 	a.value = &stringsValue{inner: value}
// 	return value
// }

// func (a *optionalArg) StringMap(target **map[string]string) *StringMap {
// 	value := &StringMap{target: target}
// 	a.value = &stringMapValue{inner: value}
// 	return value
// }
