package events

import "sync"

type Event interface{}

type Subscriber func(Event)

type Bus struct {
	mu   sync.RWMutex
	subs []Subscriber
}

func NewBus() *Bus {
	return &Bus{}
}

// Subscribe adds a subscriber callback
func (b *Bus) Subscribe(sub Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs = append(b.subs, sub)
}

// Publish sends an event to all subscribers
func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, sub := range b.subs {
		go sub(event) // async delivery (don’t block publisher)
	}
}
