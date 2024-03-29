package server

import (
	"net/http"
	"notification-service/internal/controllers"

	"github.com/gorilla/mux"
)

func addRoutes(mux *mux.Router, receiverController *controllers.ReceiverController) {
	mux.HandleFunc("/lime/eth", GetTransactionsHandler(receiverController)).Methods("GET")
}

var GetTransactionsHandler = func(receiverController *controllers.ReceiverController) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// hashes := r.URL.Query()["transactionHashes"]

			// log.Printf("%v", hashes)
			// body, err := receiverController.GetTransactions(hashes)
			// if err != nil {
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	json.NewEncoder(w).Encode(err)

			// 	return
			// }

			// w.WriteHeader(http.StatusOK)
			// json.NewEncoder(w).Encode(body)
		},
	)
}
