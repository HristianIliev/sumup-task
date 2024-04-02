package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"notification-service/internal/models/enums"
	"notification-service/internal/service"
	"notification-service/pkg/models/apimodels"
)

type NotificationController struct {
	receiverService     *service.ReceiverService
	notificationService *service.NotificationService
}

func NewNotificationController(receiverService *service.ReceiverService, notificationService *service.NotificationService) *NotificationController {
	return &NotificationController{
		receiverService:     receiverService,
		notificationService: notificationService,
	}
}

func (n *NotificationController) SendNotification(notification *apimodels.Notification) error {
	receiver, err := n.receiverService.GetReceiver(notification.ReceiverID)
	if errors.Is(service.ErrReceiverNotFound, err) {
		fmt.Println("Receiver not found")
		return err
	} else if err != nil {
		fmt.Println("other err", err)
		return err
	}

	message := &apimodels.Message{
		Content:  notification.Message,
		Email:    receiver.Email,
		Phone:    receiver.Phone,
		Endpoint: receiver.Endpoint,
		Channels: enums.Channels,
	}
	for _, preference := range receiver.ChannelPreferences {
		if preference.Preference == enums.PreferenceNever {
			message.FilterOutChannel(preference.Channel)
		}
	}

	notifications := []string{}
	for _, channel := range message.Channels {
		clone := message.Clone()
		clone.Channels = []enums.Channel{channel}

		jsonData, err := json.Marshal(clone)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		notifications = append(notifications, string(jsonData))
	}

	return n.notificationService.SendMessages(notifications)
}
