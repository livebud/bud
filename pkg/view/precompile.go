package view

// Optional method for precompiling views
type Precompiler interface {
	Precompile() (Viewer, error)
}

// Precompile the views
func Precompile(viewer Viewer) (Viewer, error) {
	if precompiler, ok := viewer.(Precompiler); ok {
		return precompiler.Precompile()
	}
	return viewer, nil
}
