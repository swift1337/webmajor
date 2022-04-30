package bus

import (
	"sync"

	"github.com/google/uuid"
)

type Bus struct {
	mu               sync.RWMutex
	subscribers      map[string]*Subscriber
	subscriberBuffer int
}

type Subscriber struct {
	id string
	ch chan any
}

func New(subscriberBuffer int) *Bus {
	return &Bus{
		mu:               sync.RWMutex{},
		subscribers:      make(map[string]*Subscriber),
		subscriberBuffer: subscriberBuffer,
	}
}

func (b *Bus) Subscribe() *Subscriber {
	sub := &Subscriber{
		id: uuid.New().String(),
		ch: make(chan any, b.subscriberBuffer),
	}

	b.mu.Lock()
	b.subscribers[sub.id] = sub
	b.mu.Unlock()

	return sub
}

func (b *Bus) Unsubscribe(sub *Subscriber) {
	close(sub.ch)

	b.mu.Lock()
	delete(b.subscribers, sub.id)
	b.mu.Unlock()
}

func (b *Bus) Publish(message any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, sub := range b.subscribers {
		select {
		case sub.ch <- message:
		default: // safe skip if subscriber buffer is ful
		}
	}
}

func (s *Subscriber) Channel() <-chan any {
	return s.ch
}
