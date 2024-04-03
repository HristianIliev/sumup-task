package sqsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"notification-service/pkg/models/apimodels"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MessagePoller interface {
	DeleteMessages([]types.Message) error
	ReceiveMessages() ([]types.Message, error)
}

type NotificationSubscriber interface {
	Should(notification *apimodels.Message) bool
	Handle(notification *apimodels.Message) error
}

type Consumer struct {
	id          int
	pollDelay   time.Duration
	sqsClient   MessagePoller
	subscribers []NotificationSubscriber
}

func NewConsumer(id int, sqsClient MessagePoller, pollDelayInSec int, subscribers []NotificationSubscriber) *Consumer {
	return &Consumer{
		id:          id,
		pollDelay:   time.Duration(pollDelayInSec) * time.Second,
		sqsClient:   sqsClient,
		subscribers: subscribers,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	log.Printf("Worker %d started listening for notifications\n", c.id)

	delayCh := time.NewTicker(c.pollDelay)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping consumer: %v\n", ctx.Err())

			delayCh.Stop()

			return
		case <-delayCh.C:
			log.Printf("Checking SQS for messages %d\n", c.id)
		}

		msgs, err := c.sqsClient.ReceiveMessages()
		if err != nil {
			log.Println("Failed to receive message: ", err)
			continue
		}

		var processedMsgs []types.Message
		for _, msg := range msgs {
			err = c.process(msg)
			if err != nil {
				log.Println("Failed to process notification. Message will be processed upon next retrieval: ", err)
				continue
			}

			processedMsgs = append(processedMsgs, msg)
		}

		err = c.sqsClient.DeleteMessages(processedMsgs)
		if err != nil {
			log.Println("Failed to delete messages: ", err)

			continue
		}
	}
}

type Notification struct {
	Message string `json:"Message"`
}

func (c *Consumer) process(message types.Message) error {
	var notification Notification
	if err := json.Unmarshal([]byte(*message.Body), &notification); err != nil {
		return fmt.Errorf("failed to decode SNS payload: %v", err)
	}

	var msg apimodels.Message
	if err := json.Unmarshal([]byte(notification.Message), &msg); err != nil {
		return fmt.Errorf("failed to decode SNS message: %v", err)
	}

	for _, subscriber := range c.subscribers {
		if subscriber.Should(&msg) {
			err := subscriber.Handle(&msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
