package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RegisterCustomerApi(router *mux.Router) {
	router.HandleFunc("/customers/{id}", customersHandler)
}

func customersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	log.Printf("CustomersHandler executed: %v", vars)
}
