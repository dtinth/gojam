package main

import "sync"

type eventBroadcaster struct {
	mutex       sync.RWMutex
	connections map[string]func(message string)

	GetWelcomeMessage func() string
}

func newEventBroadcaster() *eventBroadcaster {
	return &eventBroadcaster{
		connections: make(map[string]func(message string)),
	}
}

func (b *eventBroadcaster) Register(id string, onMessage func(message string)) func() {
	if b.GetWelcomeMessage != nil {
		onMessage(b.GetWelcomeMessage())
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.connections[id] = onMessage
	return func() {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		delete(b.connections, id)
	}
}

func (b *eventBroadcaster) Broadcast(message string) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, onMessage := range b.connections {
		onMessage(message)
	}
}
