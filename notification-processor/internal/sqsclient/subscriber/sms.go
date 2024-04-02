package subscriber

import (
	"fmt"
	"notification-service/pkg/models/apimodels"
)

type SMSSubscriber struct {
}

func (e *SMSSubscriber) Should(notification *apimodels.Message) bool {
	for _, channel := range notification.Channels {
		if channel == "sms" {
			return true
		}
	}

	return false
}

func (e *SMSSubscriber) Handle(notification *apimodels.Message) error {
	fmt.Printf("Sending sms to receiver %s with content %s\n", notification.Phone, notification.Content)

	return nil
}
