package api

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"poc/internal/model"
	"time"
)

func RegisterAppApi(router *mux.Router) {
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/health", healthHandler)

	handlers.CORS()
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
	w.WriteHeader(http.StatusOK)
	responseData := model.HttpResponse{
		Success: true,
		Timestamp: time.Now(),
		Message: "alive",
		Data:    "server01",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseData)
}
