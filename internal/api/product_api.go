package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RegisterProductApi(router *mux.Router) {
	router.HandleFunc("/products/{id}", productsHandler)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	log.Printf("ProductsHandler executed: %v", vars)
}