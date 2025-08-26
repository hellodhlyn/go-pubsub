package pubsub

import "context"

// Queue is an interface that abstracts message queue service like Amazon SQS,
// Google Cloud Pubsub, or others.
type Queue interface {
	PutMessages([]string) error
	PutMessagesWithContext(context.Context, []string) error
	FetchMessages() ([]string, error)
	FetchMessagesWithContext(ctx context.Context) ([]string, error)
}
