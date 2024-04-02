package server

import (
	"net/http"

	"notification-service/internal/controllers"

	"github.com/gorilla/mux"
)

func NewServer(receiverController *controllers.ReceiverController, notificationController *controllers.NotificationController) http.Handler {
	mux := mux.NewRouter()

	addRoutes(mux, receiverController, notificationController)

	return mux
}
