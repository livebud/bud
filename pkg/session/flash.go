package session

type Flash struct {
	wrapper *wrapper
}

func (f *Flash) Add(message string) {
	f.wrapper.Flashes = append(f.wrapper.Flashes, message)
}
