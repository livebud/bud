package commander

type Args struct {
	Name  string
	value value
}

func (a *Args) Strings(target *[]string) *Strings {
	*target = []string{}
	value := &Strings{target: target}
	a.value = &stringsValue{inner: value}
	return value
}
