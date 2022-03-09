package commander

type Args struct {
	Name  string
	value value
}

func (a *Args) Strings(target *[]string) *Strings {
	value := &Strings{target: target}
	a.value = &stringsValue{inner: value}
	return value
}
