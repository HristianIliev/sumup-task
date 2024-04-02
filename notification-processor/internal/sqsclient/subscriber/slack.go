package subscriber

import (
	"fmt"
	"notification-service/pkg/models/apimodels"
)

type SlackSubscriber struct {
}

func (e *SlackSubscriber) Should(notification *apimodels.Message) bool {
	for _, channel := range notification.Channels {
		if channel == "slack" {
			return true
		}
	}

	return false
}

func (e *SlackSubscriber) Handle(notification *apimodels.Message) error {
	fmt.Printf("Sending slack message to receiver %s with content %s\n", notification.Endpoint, notification.Content)

	return nil
}
