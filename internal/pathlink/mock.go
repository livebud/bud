package pathlink

func Mock(fallback Linker) *Mocker {
	return &Mocker{linker: fallback}
}

type Mocker struct {
	linker   Linker
	MockLink func(from string, to ...string) error
}

var _ Linker = (*Mocker)(nil)

func (m *Mocker) Link(from string, to ...string) error {
	if m.MockLink != nil {
		return m.MockLink(from, to...)
	}
	return m.linker.Link(from, to...)
}
