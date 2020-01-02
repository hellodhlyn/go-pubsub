package pubsub

// Queue is an interface that abstracts message queue service like Amazon SQS,
// Google Cloud Pubsub, or others.
type Queue interface {
	PutMessages([]string) error
	FetchMessages() ([]string, error)
}
