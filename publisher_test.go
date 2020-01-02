package pubsub_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hellodhlyn/go-pubsub"
)

func TestPublisher_Publish_Success(t *testing.T) {
	err := pubsub.NewPublisher(validQueue{}).Publish("dummy-message")
	assert.Nil(t, err)
}

func TestPublisher_Publish_Error(t *testing.T) {
	err := pubsub.NewPublisher(errorQueue{}).Publish("dummy-message")
	assert.EqualError(t, err, "something went wrong :(")
}

func TestPublisher_PublishMany_Success(t *testing.T) {
	err := pubsub.NewPublisher(validQueue{}).PublishMany([]string{"dummy-message", "dummy-message", "dummy-message"})
	assert.Nil(t, err)
}

func TestPublisher_PublishMany_Error(t *testing.T) {
	err := pubsub.NewPublisher(errorQueue{}).PublishMany([]string{"dummy-message", "dummy-message", "dummy-message"})
	assert.EqualError(t, err, "something went wrong :(")
}
