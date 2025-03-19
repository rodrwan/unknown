package pubsub

import (
	"cryptodashboard/internal/services"
	"sync"
	"time"
)

type Subscriber struct {
	ID       string
	Channel  chan []*services.Balance
	Callback func([]*services.Balance)
}

type PubSub struct {
	subscribers map[string][]*Subscriber
	mu          sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]*Subscriber),
	}
}

func (ps *PubSub) Subscribe(topic string, callback func([]*services.Balance)) *Subscriber {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := &Subscriber{
		ID:       generateID(),
		Channel:  make(chan []*services.Balance, 1),
		Callback: callback,
	}

	ps.subscribers[topic] = append(ps.subscribers[topic], sub)
	return sub
}

func (ps *PubSub) Unsubscribe(topic string, sub *Subscriber) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs := ps.subscribers[topic]
	for i, s := range subs {
		if s.ID == sub.ID {
			ps.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			close(s.Channel)
			return
		}
	}
}

func (ps *PubSub) Publish(topic string, data []*services.Balance) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subs := ps.subscribers[topic]
	for _, sub := range subs {
		go func(s *Subscriber) {
			s.Callback(data)
		}(sub)
	}
}

func generateID() string {
	return "sub-" + time.Now().Format("20060102150405")
}
