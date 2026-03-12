package sse

import "sync"

// Event is the JSON payload sent to SSE clients.
type Event struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

// Broker manages per-user SSE subscriber channels.
type Broker struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe registers a new channel for the given user and returns it.
func (b *Broker) Subscribe(userID string) chan Event {
	ch := make(chan Event, 4)
	b.mu.Lock()
	b.subscribers[userID] = append(b.subscribers[userID], ch)
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes the channel from the user's subscribers and closes it.
func (b *Broker) Unsubscribe(userID string, ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	subs := b.subscribers[userID]
	for i, s := range subs {
		if s == ch {
			b.subscribers[userID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			return
		}
	}
}

// Publish sends an event to all subscribers of the given user (non-blocking).
func (b *Broker) Publish(userID string, event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subscribers[userID] {
		select {
		case ch <- event:
		default:
		}
	}
}
