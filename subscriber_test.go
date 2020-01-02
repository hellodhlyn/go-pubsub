package pubsub_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hellodhlyn/go-pubsub"
)

func TestSubscriber_SubscribeFunc_Success(t *testing.T) {
	subscriber := pubsub.NewSubscriber(validQueue{})

	// Send SIGTERM after 10 ms.
	pid := os.Getpid()
	go func() { time.Sleep(10 * time.Millisecond); _ = syscall.Kill(pid, syscall.SIGTERM) }()

	msgCnt := 0
	fn := func(message string) {
		msgCnt += 1
		assert.Equal(t, message, "dummy-message")
	}
	err := subscriber.SubscribeFunc(fn)

	assert.EqualError(t, err, "received signal: terminated")
	assert.NotEqual(t, msgCnt, 0)
}

func TestSubscriber_SubscribeFunc_Error(t *testing.T) {
	subscriber := pubsub.NewSubscriber(errorQueue{})

	msgCnt := 0
	err := subscriber.SubscribeFunc(func(_ string) { msgCnt += 1 })

	assert.EqualError(t, err, "something went wrong :(")
	assert.Equal(t, msgCnt, 0)
}
