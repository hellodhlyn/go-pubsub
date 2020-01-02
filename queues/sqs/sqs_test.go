package sqs_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hellodhlyn/go-pubsub/queues/sqs"
)

var (
	queueURL = "https://dummy-queue-url"
	msgCount = 0
)

//// Mocking SQS Queue
type mockSQS struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (m mockSQS) GetQueueUrl(_ *awssqs.GetQueueUrlInput) (*awssqs.GetQueueUrlOutput, error) {
	return &awssqs.GetQueueUrlOutput{QueueUrl: aws.String(queueURL)}, nil
}

func (m mockSQS) SendMessageBatch(input *awssqs.SendMessageBatchInput) (*awssqs.SendMessageBatchOutput, error) {
	m.Called(input)
	msgCount += len(input.Entries)
	return &awssqs.SendMessageBatchOutput{}, nil
}

func (m mockSQS) ReceiveMessage(input *awssqs.ReceiveMessageInput) (*awssqs.ReceiveMessageOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*awssqs.ReceiveMessageOutput), nil
}

func (m mockSQS) DeleteMessageBatch(input *awssqs.DeleteMessageBatchInput) (*awssqs.DeleteMessageBatchOutput, error) {
	m.Called(input)
	return &awssqs.DeleteMessageBatchOutput{}, nil
}

func newMockQueue(fn func(*mockSQS)) (*sqs.SQSQueue, func(*testing.T)) {
	var mockClient = mockSQS{}
	fn(&mockClient)
	return &sqs.SQSQueue{URL: queueURL, Client: mockClient}, func(t *testing.T) { mockClient.AssertExpectations(t) }
}

/// Tests
func TestSQSQueue_PutMessages_Success1(t *testing.T) {
	mockedQueue, assertStubs := newMockQueue(func(client *mockSQS) {
		client.On("SendMessageBatch", mock.Anything).Once()
		// FIXME - id of messages are generated uuid so we can't expect the value
	})

	err := mockedQueue.PutMessages(make([]string, 3))

	assert.Nil(t, err)
	assert.Equal(t, msgCount, 3)
	assertStubs(t)
}

// If the size of messages is larger than 10, the messages have been split by
// 10 and call API multiple times.
func TestSQSQueue_PutMessages_Success2(t *testing.T) {
	msgCount = 0
	mockedQueue, assertStubs := newMockQueue(func(client *mockSQS) {
		client.On("SendMessageBatch", mock.Anything).Times(2)
		// FIXME - id of messages are generated uuid so we can't expect the value
	})

	err := mockedQueue.PutMessages(make([]string, 15))

	assert.Nil(t, err)
	assert.Equal(t, msgCount, 15)
	assertStubs(t)
}

func TestSQSQueue_FetchMessages_Success1(t *testing.T) {
	mockQueue, assertStubs := newMockQueue(func(client *mockSQS) {
		expect1 := &awssqs.ReceiveMessageInput{MaxNumberOfMessages: aws.Int64(10), QueueUrl: &queueURL, WaitTimeSeconds: aws.Int64(20)}
		client.On("ReceiveMessage", expect1).Return(makeReceiveMessageOutput(3), nil).Once()

		expect2 := &awssqs.DeleteMessageBatchInput{
			Entries:  []*awssqs.DeleteMessageBatchRequestEntry{{Id: aws.String("dummy-id")}, {Id: aws.String("dummy-id")}, {Id: aws.String("dummy-id")}},
			QueueUrl: &queueURL,
		}
		client.On("DeleteMessageBatch", expect2).Once()
	})

	messages, err := mockQueue.FetchMessages()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(messages), "number of messages should be 3")
	assertStubs(t)
}

func TestSQSQueue_FetchMessages_Success2(t *testing.T) {
	mockQueue, assertStubs := newMockQueue(func(client *mockSQS) {
		expect := &awssqs.ReceiveMessageInput{MaxNumberOfMessages: aws.Int64(10), QueueUrl: &queueURL, WaitTimeSeconds: aws.Int64(20)}
		client.On("ReceiveMessage", expect).Return(makeReceiveMessageOutput(0), nil).Once()
	})

	messages, err := mockQueue.FetchMessages()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messages), "number of messages should be 0")
	assertStubs(t)
}

func makeReceiveMessageOutput(count int) *awssqs.ReceiveMessageOutput {
	messages := make([]*awssqs.Message, count)
	for i := 0; i < len(messages); i += 1 {
		messages[i] = &awssqs.Message{MessageId: aws.String("dummy-id"), Body: aws.String("dummy-received-message")}
	}

	return &awssqs.ReceiveMessageOutput{Messages: messages}
}
