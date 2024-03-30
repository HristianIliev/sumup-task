package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notification-service/internal/controllers"
	"notification-service/internal/models/apimodels"

	"github.com/gorilla/mux"
)

func addRoutes(mux *mux.Router, receiverController *controllers.ReceiverController) {
	mux.HandleFunc("/receivers", CreateReceiverHandler(receiverController)).Methods("POST")
	mux.HandleFunc("/receivers/{id}", GetReceiverHandler(receiverController)).Methods("GET")
}

var CreateReceiverHandler = func(receiverController *controllers.ReceiverController) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var receiver apimodels.Receiver
			err := json.NewDecoder(r.Body).Decode(&receiver)
			if err != nil {
				fmt.Println("here")
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
