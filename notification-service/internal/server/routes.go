package server

import (
	"encoding/json"
	"net/http"
	"notification-service/internal/controllers"
	"notification-service/pkg/models/apimodels"

	"github.com/gorilla/mux"
)

func addRoutes(mux *mux.Router, receiverController *controllers.ReceiverController, notificationController *controllers.NotificationController) {
	mux.HandleFunc("/receivers", CreateReceiverHandler(receiverController)).Methods("POST")
	mux.HandleFunc("/receivers/{id}", GetReceiverHandler(receiverController)).Methods("GET")
	mux.HandleFunc("/notifications", SendNotificationHandler(notificationController)).Methods("POST")
}

var CreateReceiverHandler = func(receiverController *controllers.ReceiverController) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var receiver apimodels.Receiver
			err := json.NewDecoder(r.Body).Decode(&receiver)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)

				return
			}

			body, err := receiverController.RegisterReceiver(&receiver)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)

				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(body)
		},
	)
}

var GetReceiverHandler = func(receiverController *controllers.ReceiverController) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			response, err := receiverController.GetReceiver(vars["id"])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)

				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		},
	)
}

var SendNotificationHandler = func(notificationController *controllers.NotificationController) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var notification apimodels.Notification
			err := json.NewDecoder(r.Body).Decode(&notification)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)

				return
			}

			err = notificationController.SendNotification(&notification)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)

				return
			}

			w.WriteHeader(http.StatusOK)
		},
	)
}
