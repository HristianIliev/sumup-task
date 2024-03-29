package server

import (
	"encoding/json"
	"ethereum-fetcher/internal/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func addRoutes(mux *mux.Router, transactionController *controllers.TransactionController) {
	mux.HandleFunc("/lime/eth", GetTransactionsHandler(transactionController)).Methods("GET")
}

var GetTransactionsHandler = func(transactionController *controllers.TransactionController) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			hashes := r.URL.Query()["transactionHashes"]

			log.Printf("%v", hashes)
			body, err := transactionController.GetTransactions(hashes)
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
