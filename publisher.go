package pubsub

// Publisher puts messages into the message queue.
type Publisher struct {
	queue Queue
}

// NewPublisher creates a new Publisher using the given queue.
func NewPublisher(queue Queue) *Publisher {
	return &Publisher{queue}
}

// Publish puts the message into the queue.
func (p *Publisher) Publish(message string) error {
	return p.queue.PutMessages([]string{message})
}

// PublishMany puts multiple messages into the queue.
func (p *Publisher) PublishMany(messages []string) error {
	return p.queue.PutMessages(messages)
}
