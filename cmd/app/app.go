package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"poc/internal/api"
	"time"
)

func init() {
	log.Println("Initializing the application...")
	initAppConfig()
}

func main() {
	log.Println("Building application server...")
	buildAndRunServer()
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

func buildAndRunServer() {
	router := mux.NewRouter()
	registerApis(router)

	originsOk := handlers.AllowedOrigins([]string{"*"})
	handlers.CORS(originsOk)(router)

	requestHandler := handlers.LoggingHandler(os.Stdout, router)
	serverPort :=  viper.GetString("Port")

	log.Printf("Starting server on port %v", serverPort)
	srv := &http.Server{
		Handler:      requestHandler,
		Addr:         fmt.Sprintf(":%v", serverPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func registerApis(router *mux.Router) {
	api.RegisterAppApi(router)
	api.RegisterAuthApi(router)
	api.RegisterCustomerApi(router)
	api.RegisterProductApi(router)
}