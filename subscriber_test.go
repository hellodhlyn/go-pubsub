package pubsub_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hellodhlyn/go-pubsub"
)

func TestSubscriber_SubscribeFunc_Success(t *testing.T) {
	subscriber := pubsub.NewSubscriber(validQueue{})

	msgCnt := 0

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	fn := func(message string) {
		msgCnt += 1
		assert.Equal(t, message, "dummy-message")

		if msgCnt == 3 {
			cancel()
			done <- struct{}{}
		}
	}

	go func() {
		err := subscriber.SubscribeFunc(ctx, fn)
		assert.ErrorIs(t, err, context.Canceled)
	}()

	<-done
}

func TestSubscriber_SubscribeFunc_Error(t *testing.T) {
	subscriber := pubsub.NewSubscriber(errorQueue{})

	msgCnt := 0
	err := subscriber.SubscribeFunc(context.Background(), func(_ string) { msgCnt += 1 })

	assert.EqualError(t, err, "something went wrong :(")
	assert.Equal(t, msgCnt, 0)
}
