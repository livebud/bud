package mod

// New modfile loader
func New() *Finder {
	return &Finder{}
}

// Finder struct
type Finder struct {
}

// Find a modfile
func (f *Finder) Find(dir string) (*File, error) {
	modfile, err := FindIn(dir)
	if err != nil {
		return nil, err
	}
	return modfile, nil
}
