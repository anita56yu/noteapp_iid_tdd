package usecase

import (
	"sync"
)

type Event any

type Subscriber[T Event] func(event T)

type EventBus[T Event] interface {
	Publish(event T)
	Subscribe(handler Subscriber[T])
}

type InMemoryEventBus[T Event] struct {
	mu          sync.RWMutex
	subscribers []Subscriber[T]
}

func NewEventBus[T Event]() *InMemoryEventBus[T] {
	return &InMemoryEventBus[T]{
		subscribers: make([]Subscriber[T], 0),
	}
}

func (eb *InMemoryEventBus[T]) Subscribe(handler Subscriber[T]) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers = append(eb.subscribers, handler)
}

func (eb *InMemoryEventBus[T]) Publish(event T) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	for _, handler := range eb.subscribers {
		go handler(event)
	}
}
