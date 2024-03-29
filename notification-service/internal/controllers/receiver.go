package controllers

import "notification-service/internal/service/preferences"

type ReceiverController struct {
}

func New(preferencesService *preferences.PreferencesService) *ReceiverController {
	return &ReceiverController{}
}

func (r *ReceiverController) RegisterReceiver() {
}
