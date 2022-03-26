package pubsub

import (
	"sync"
)

type Publisher interface {
	Publish(topic string, data []byte)
}

type Subscriber interface {
	Subscribe(topics ...string) Subscription
}

type Client interface {
	Publisher
	Subscriber
}

type Subscription interface {
	Wait() <-chan []byte
	Close()
}

func New() *Memory {
	return &Memory{
		topics: map[string]map[int]chan []byte{},
	}
}

type Memory struct {
	mu     sync.RWMutex
	topics map[string]map[int]chan []byte
	cid    int
}

var _ Publisher = (*Memory)(nil)
var _ Subscriber = (*Memory)(nil)

func (m *Memory) Publish(topic string, data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	subs := m.topics[topic]
	for _, ch := range subs {
		select {
		case ch <- data:
		default:
			// Disregard slow subscribers
		}
	}
}

func (m *Memory) Subscribe(topics ...string) Subscription {
	ch := make(chan []byte, 1)
	m.mu.Lock()
	defer m.mu.Unlock()
	cid := m.cid
	for _, topic := range topics {
		if m.topics[topic] == nil {
			m.topics[topic] = map[int]chan []byte{}
		}
		m.topics[topic][cid] = ch
	}
	m.cid++
	return &subscriber{ch, m.close(cid, topics)}
}

func (m *Memory) close(MemoryID int, topics []string) func() {
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		for i, topic := range topics {
			// Same channel across subcribers, so we should only close once
			if i == 0 {
				ch := m.topics[topic][MemoryID]
				close(ch)
			}
			delete(m.topics[topic], MemoryID)
		}
	}
}

type subscriber struct {
	ch     chan []byte
	closer func()
}

func (s *subscriber) Wait() <-chan []byte {
	return s.ch
}

func (s *subscriber) Close() {
	s.closer()
}
