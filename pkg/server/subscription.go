package server

import "github.com/cigarframework/reconciled/pkg/api"

type subscriptionManager struct {
	input chan *api.Notification
	reg   chan chan<- *api.Notification
	unreg chan chan<- *api.Notification

	outputs map[chan<- *api.Notification]bool
}

type SubscriptionManager interface {
	Subscribe(chan<- *api.Notification)
	Cancel(chan<- *api.Notification)
	Close() error
	Publish(*api.Notification)
}

func (b *subscriptionManager) broadcast(m *api.Notification) {
	for ch := range b.outputs {
		ch <- m
	}
}

func (b *subscriptionManager) run() {
	for {
		select {
		case m := <-b.input:
			b.broadcast(m)
		case ch, ok := <-b.reg:
			if ok {
				b.outputs[ch] = true
			} else {
				return
			}
		case ch := <-b.unreg:
			delete(b.outputs, ch)
		}
	}
}

func newSubscriptionManager(buflen int) SubscriptionManager {
	b := &subscriptionManager{
		input:   make(chan *api.Notification, buflen),
		reg:     make(chan chan<- *api.Notification),
		unreg:   make(chan chan<- *api.Notification),
		outputs: make(map[chan<- *api.Notification]bool),
	}

	go b.run()
	return b
}

func (b *subscriptionManager) Subscribe(ch chan<- *api.Notification) {
	b.reg <- ch
}

func (b *subscriptionManager) Cancel(ch chan<- *api.Notification) {
	b.unreg <- ch
}

func (b *subscriptionManager) Close() error {
	close(b.reg)
	return nil
}

func (b *subscriptionManager) Publish(m *api.Notification) {
	if b != nil {
		b.input <- m
	}
}
