package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/google/uuid"
)

// SQSQueue is a Queue with Amazon SQS.
type SQSQueue struct {
	URL    string
	Client sqsiface.SQSAPI
}

type SQSQueueConfig struct {
	QueueName  string
	Region     string
	MaxRetries int
}

// NewSQSQueue creates a Queue with Amazon SQS.
func NewSQSQueue(cfg *SQSQueueConfig) (*SQSQueue, error) {
	sess, err := session.NewSession(&aws.Config{Region: &cfg.Region, MaxRetries: &cfg.MaxRetries})
	if err != nil {
		return nil, err
	}
	client := sqs.New(sess)

	output, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: &cfg.QueueName})
	if err != nil {
		return nil, err
	}

	return &SQSQueue{URL: *output.QueueUrl, Client: client}, nil
}

// PutMessages enqueue given messages to SQS queue.
//
// Since SQS supports at maximum 10 messages for one batch operation, we split
// the given slice by 10 times.
func (q *SQSQueue) PutMessages(messages []string) error {
	for i := 0; i <= (len(messages)-1)/10; i += 1 {
		fromIdx := i * 10
		toIdx := fromIdx + 10
		if toIdx > len(messages) {
			toIdx = len(messages)
		}

		if err := q.putMessages(context.Background(), messages[fromIdx:toIdx]); err != nil {
			return err
		}
	}
	return nil
}

// PutMessagesWithContext enqueue given messages to SQS queue.
//
// Since SQS supports at maximum 10 messages for one batch operation, we split
// the given slice by 10 times.
func (q *SQSQueue) PutMessagesWithContext(ctx context.Context, messages []string) error {
	for i := 0; i <= (len(messages)-1)/10; i += 1 {
		fromIdx := i * 10
		toIdx := fromIdx + 10
		if toIdx > len(messages) {
			toIdx = len(messages)
		}

		if err := q.putMessages(ctx, messages[fromIdx:toIdx]); err != nil {
			return err
		}
	}
	return nil
}

func (q *SQSQueue) putMessages(ctx context.Context, messages []string) error {
	entries := make([]*sqs.SendMessageBatchRequestEntry, len(messages))
	for i, message := range messages {
		entries[i] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(uuid.New().String()),
			MessageBody: aws.String(message),
		}
	}

	_, err := q.Client.SendMessageBatchWithContext(ctx, &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: &q.URL,
	})
	return err
}

// FetchMessages dequeue messages from SQS queue using long polling.
func (q *SQSQueue) FetchMessages() ([]string, error) {
	return q.fetchMessages(context.Background())
}

// FetchMessagesWithContext dequeue messages from SQS queue using long polling.
func (q *SQSQueue) FetchMessagesWithContext(ctx context.Context) ([]string, error) {
	return q.fetchMessages(ctx)
}

func (q *SQSQueue) fetchMessages(ctx context.Context) ([]string, error) {
	output, err := q.Client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(10),
		QueueUrl:            &q.URL,
		WaitTimeSeconds:     aws.Int64(20),
	})
	if err != nil {
		return nil, err
	}
	if len(output.Messages) == 0 {
		return []string{}, nil
	}

	messages := make([]string, len(output.Messages))
	entries := make([]*sqs.DeleteMessageBatchRequestEntry, len(output.Messages))
	for i, m := range output.Messages {
		messages[i] = *m.Body
		entries[i] = &sqs.DeleteMessageBatchRequestEntry{Id: m.MessageId, ReceiptHandle: m.ReceiptHandle}
	}

	_, err = q.Client.DeleteMessageBatchWithContext(ctx, &sqs.DeleteMessageBatchInput{Entries: entries, QueueUrl: &q.URL})
	if err != nil {
		return nil, err
	}
	return messages, nil
}
