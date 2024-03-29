package server

import (
	"net/http"

	"ethereum-fetcher/internal/controllers"

	"github.com/gorilla/mux"
)

func NewServer(transactionController *controllers.TransactionController) http.Handler {
	mux := mux.NewRouter()

	addRoutes(mux, transactionController)

	return mux
}
