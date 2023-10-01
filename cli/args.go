package cli

type Args struct {
	Name  string
	value value
}

func (a *Args) Optional() *OptionalArgs {
	return &OptionalArgs{a}
}

func (a *Args) Strings(target *[]string) *Strings {
	*target = []string{}
	value := &Strings{target: target}
	a.value = &stringsValue{inner: value}
	return value
}

func (a *Args) verify(name string) error {
	return a.value.verify("<" + name + "...>")
}

type OptionalArgs struct {
	a *Args
}

func (a *OptionalArgs) Strings(target *[]string) *Strings {
	*target = []string{}
	value := &Strings{target: target, optional: true}
	a.a.value = &stringsValue{inner: value}
	return value
}
