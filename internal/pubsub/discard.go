package pubsub

func Discard() *discarder {
	return &discarder{}
}

type discarder struct {
}

var _ Publisher = (*discarder)(nil)
var _ Subscriber = (*discarder)(nil)

func (d *discarder) Publish(topic string, data []byte) {
}

func (d *discarder) Subscribe(topics ...string) Subscription {
	ch := make(chan []byte)
	close(ch) // Close the channel immediately
	closer := func() {}
	return &discardable{ch, closer}
}

type discardable struct {
	ch     chan []byte
	closer func()
}

var _ Subscription = (*discardable)(nil)

func (s *discardable) Wait() <-chan []byte {
	return s.ch
}

func (s *discardable) Close() {
	s.closer()
}
