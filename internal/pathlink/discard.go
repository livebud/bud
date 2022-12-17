package pathlink

var Discard Linker = discard{}

type discard struct{}

func (discard) Link(from string, to ...string) error {
	return nil
}
