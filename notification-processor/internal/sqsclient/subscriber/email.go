package subscriber

import (
	"fmt"
	"notification-service/pkg/models/apimodels"
)

type EmailSubscriber struct {
}

func (e *EmailSubscriber) Should(notification *apimodels.Message) bool {
	for _, channel := range notification.Channels {
		if channel == "email" {
			return true
		}
	}

	return false
}

func (e *EmailSubscriber) Handle(notification *apimodels.Message) error {
	fmt.Printf("Sending email to receiver %s with content %s\n", notification.Email, notification.Content)

	return nil
}
