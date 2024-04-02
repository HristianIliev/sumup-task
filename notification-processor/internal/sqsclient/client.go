package sqsclient

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	queueURL string
	client   *sqs.Client
}

func New(client *sqs.Client, queueURL string) *SQSClient {
	return &SQSClient{
		queueURL: queueURL,
		client:   client,
	}
}

func (c *SQSClient) ReceiveMessages() ([]types.Message, error) {
	ctx := context.Background()
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              &c.queueURL,
		MaxNumberOfMessages:   10,
		MessageAttributeNames: []string{"All"},
	})
	if err != nil {
		return nil, err
	}

	if len(output.Messages) == 0 {
		fmt.Println("No new SQS Messages")

		return nil, nil
	}

	fmt.Printf("Received %d messages.\n", len(output.Messages))

	return output.Messages, nil
}

func (c *SQSClient) DeleteMessages(msgs []types.Message) error {
	for _, msg := range msgs {
		ctx := context.Background()
		timeout := time.Second * 10
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &c.queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			return err
		}

		fmt.Printf("SQS message with id %s was deleted\n", *msg.MessageId)
	}

	return nil
}
