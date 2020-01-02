package pubsub

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
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
// function. This is blocking action; it unblocks if either of SIGINT,
// SIGTERM, SIGKILL has received.
func (s *Subscriber) SubscribeFunc(fn func(string)) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	for {
		messages, err := s.queue.FetchMessages()
		if err != nil {
			return err
		}

		for _, msg := range messages {
			fn(msg)
		}

		select {
		case sig := <-sigCh:
			return errors.New("received signal: " + sig.String())
		default:
			continue
		}
	}
}
