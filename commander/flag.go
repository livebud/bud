package commander

import "flag"

type Flag struct {
	name  string
	usage string
	value flag.Getter
}

func (f *Flag) Int(target *int) *Int {
	value := &Int{target: target}
	f.value = &intValue{value}
	return value
}

func (f *Flag) String(target *string) *String {
	value := &String{target: target}
	f.value = &stringValue{inner: value}
	return value
}

type value interface {
	flag.Getter
	verify(name string) error
}

func (f *Flag) verify(name string) error {
	if verifier, ok := f.value.(value); ok {
		return verifier.verify(name)
	}
	return nil
}

func verifyFlags(flags []*Flag) error {
	for _, flag := range flags {
		if err := flag.verify(flag.name); err != nil {
			return err
		}
	}
	return nil
}
