package mod

// New modfile loader
func New() *Finder {
	return &Finder{
		cache: newCache(),
	}
}

// Finder struct
type Finder struct {
	cache *cache
}

// Find a modfile
func (f *Finder) Find(dir string) (File, error) {
	if modfile, ok := f.cache.Get(dir); ok {
		return modfile, nil
	}
	modfile, err := Find(dir)
	if err != nil {
		return nil, err
	}
	f.cache.Set(dir, modfile)
	return modfile, nil
}
