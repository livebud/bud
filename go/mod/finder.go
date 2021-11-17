package mod

import "gitlab.com/mnm/bud/internal/modcache"

// New modfile loader
func New(cache *modcache.Cache) *Finder {
	return &Finder{cache}
}

// Finder struct
type Finder struct {
	cache *modcache.Cache
}

// Find a modfile
func (f *Finder) Find(dir string) (*File, error) {
	modfile, err := FindIn(f.cache, dir)
	if err != nil {
		return nil, err
	}
	return modfile, nil
}
