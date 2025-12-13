package main

import (
	"fmt"
	"net/http"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/infra/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	startServer()
}

func startServer() {
	cepHandler := api.NewCEPHandler().GetCurrentWeather
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.HandleFunc("/", cepHandler)

	fmt.Println("Starting web server on port", ":8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
