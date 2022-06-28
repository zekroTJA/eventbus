package eventbus

import (
	"sync"

	"github.com/rs/xid"
)

// EventBus is an n-to-n pub-sub structure to
// send and receive event messages between
// services.
type EventBus[T any] struct {
	buffSize int

	m    sync.RWMutex
	subs map[string]chan T
}

// New returns a new instance of EventBus with
// the given type T as message type.
//
// Optionally, you can pass a custon buffer size
// used for the event subscription channels. It
// defines the number of messages which can be
// stored in each subscription channel without
// blocking the publishing go routine until the
// message has been picked up by the subscribing
// go routine. The default value when not passed
// is 100. When the EventBus might experience
// high traffic, it is recommendet to allocate
// larger buffer sized to avoid blocking the
// publishing go routine.
func New[T any](buffSize ...int) *EventBus[T] {
	return &EventBus[T]{
		buffSize: opt(buffSize, 100),
		subs:     map[string]chan T{},
	}
}

// Publish sends the passed message to each
// subscription to the EventBus.
//
// If no one subscribed to the EventBus, the
// message will be discarded.
func (t *EventBus[T]) Publish(v T) {
	t.m.RLock()
	defer t.m.RUnlock()

	for _, s := range t.subs {
		s <- v
	}
}

// Subscribe returns a channel receiving messages
// published to the EventBus as well as a function
// to unsubscribe.
//
// When the unsubscribe function is called, the
// channel will be closed and the subscription will
// be removed from the EventBus.
func (t *EventBus[T]) Subscribe() (<-chan T, func()) {
	id := xid.New().String()

	ch := make(chan T, t.buffSize)

	t.m.Lock()
	defer t.m.Unlock()

	t.subs[id] = ch

	unsub := func() {
		delete(t.subs, id)
		close(ch)
	}

	return ch, unsub
}

// SubscribeFunc is a wrapper for Subscribe which
// calls the given function f when a message has
// been received passing the received message.
//
// The returned function unsubscribes from the
// EventBus.
func (t *EventBus[T]) SubscribeFunc(f func(T)) func() {
	if f == nil {
		panic("f must not be nil")
	}

	ch, unsub := t.Subscribe()

	go func() {
		for v := range ch {
			f(v)
		}
	}()

	return unsub
}
