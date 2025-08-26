package pubsub_test

import (
	"context"
	"errors"
)

type validQueue struct{}
type errorQueue struct{}

func (validQueue) PutMessages(_ []string) error                                 { return nil }
func (validQueue) PutMessagesWithContext(ctx context.Context, _ []string) error { return nil }
func (validQueue) FetchMessages() ([]string, error) {
	return []string{"dummy-message", "dummy-message", "dummy-message"}, nil
}
func (validQueue) FetchMessagesWithContext(ctx context.Context) ([]string, error) {
	return []string{"dummy-message", "dummy-message", "dummy-message"}, nil
}

func (errorQueue) PutMessages(_ []string) error { return errors.New("something went wrong :(") }
func (errorQueue) PutMessagesWithContext(ctx context.Context, _ []string) error {
	return errors.New("something went wrong :(")
}
func (errorQueue) FetchMessages() ([]string, error) {
	return nil, errors.New("something went wrong :(")
}
func (errorQueue) FetchMessagesWithContext(ctx context.Context) ([]string, error) {
	return nil, errors.New("something went wrong :(")
}
