package main

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"poc/internal/api"
	"time"
)

func init() {
	log.Println("Initializing the application...")
	initAppConfig()
}

func main() {
	log.Println("Building application server...")

	router := mux.NewRouter()
	registerApis(router)

	log.Println("Starting server on port 8080")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func initAppConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Application config not found: ")
		}
		panic(err)
	}

	log.Printf("ClientId: %v\n", viper.Get("ClientId"))
	log.Printf("ClientSecret: %v\n", viper.Get("ClientSecret"))
	log.Printf("Issuer: %v\n", viper.Get("Issuer"))
}

func registerApis(router *mux.Router) {
	api.RegisterAppApi(router)
	api.RegisterAuthApi(router)
	api.RegisterCustomerApi(router)
	api.RegisterProductApi(router)
}