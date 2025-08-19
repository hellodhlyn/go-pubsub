package pubsub

import (
	"context"
)

// Subscriber fetches messages from the message queue and passes it to the
// function or channel.
type Subscriber struct {
	queue Queue
}

// NewSubscriber creates a new Subscriber using the given queue.
func NewSubscriber(queue Queue) *Subscriber {
	return &Subscriber{queue}
}

// SubscribeFunc fetch messages from the queue, and pass it to the given
// function.
//
// The loop runs until an error occurs or the given context is canceled.
// If FetchMessagesWithContext returns an error, the function stops and
// returns that error immediately.
//
// The handler function is called once for each message retrieved.
func (s *Subscriber) SubscribeFunc(ctx context.Context, fn func(string)) error {
	for {
		messages, err := s.queue.FetchMessagesWithContext(ctx)
		if err != nil {
			return err
		}

		for _, msg := range messages {
			fn(msg)
		}
	}
}
