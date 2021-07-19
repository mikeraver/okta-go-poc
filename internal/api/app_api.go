package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"poc/internal/model"
)

func RegisterAppApi(router *mux.Router) {
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/health", healthHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	type customData struct {
		Profile         map[string]string
		IsAuthenticated bool
	}

	data := customData{
		Profile:         getProfileData(r),
		IsAuthenticated: isAuthenticated(r),
	}
	tpl.ExecuteTemplate(w, "home.gohtml", data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HomeHandler executed")
	w.WriteHeader(http.StatusOK)
	responseData := model.HttpResponse{
		Success: true,
		Message: "alive",
		Data:    "server01",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseData)
}
