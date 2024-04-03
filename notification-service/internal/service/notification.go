package service

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
)

type NotificationService struct {
	snsClient *sns.Client
	topic     string
}

func NewNotificationService(snsClient *sns.Client, topic string) *NotificationService {
	return &NotificationService{
		snsClient: snsClient,
		topic:     topic,
	}
}

func (n *NotificationService) SendMessages(messages []string) error {
	for _, msg := range messages {
		UUID, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		_, err = n.snsClient.Publish(context.TODO(), &sns.PublishInput{
			Message:                &msg,
			MessageGroupId:         aws.String(UUID.String()),
			MessageDeduplicationId: aws.String(UUID.String()),
			TopicArn:               &n.topic,
		})

		if err != nil {
			log.Println("Error publishing message to SNS topic:", err)
			return err
		}

		log.Printf("Successfully sent %s\n", msg)
	}

	return nil
}
