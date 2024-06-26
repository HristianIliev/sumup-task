package controllers

import (
	"notification-service/internal/adapters"
	"notification-service/internal/service"
	"notification-service/pkg/models/apimodels"
)

type ReceiverController struct {
	receiverService *service.ReceiverService
}

func NewReceiverController(receiverService *service.ReceiverService) *ReceiverController {
	return &ReceiverController{
		receiverService: receiverService,
	}
}

func (r *ReceiverController) GetReceiver(id string) (*apimodels.Receiver, error) {
	receiver, err := r.receiverService.GetReceiver(id)
	if err != nil {
		return nil, err
	}

	return adapters.DbReceiverToApiReceiver(receiver), nil
}

func (r *ReceiverController) RegisterReceiver(newReceiver *apimodels.Receiver) (*apimodels.Receiver, error) {
	receiver := adapters.ApiReceiverToDbReceiver(newReceiver)
	result, err := r.receiverService.CreateReceiver(receiver)
	if err != nil {
		return nil, err
	}

	return adapters.DbReceiverToApiReceiver(result), nil
}

func (r *ReceiverController) UpdateChannels() (*apimodels.Receiver, error) {
	return nil, nil
}
