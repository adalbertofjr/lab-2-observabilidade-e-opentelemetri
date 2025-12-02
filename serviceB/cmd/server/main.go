package main

import (
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/cmd/configs"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/api"
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/gateway"
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/web"
	usecase "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/usecase/weather"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	startServer(configs)
}

func startServer(configs *configs.Conf) {
	weatherGateway := gateway.NewWeatherAPI(configs.WeatherAPIKey)
	weatherUseCase := usecase.NewWeatherUseCase(weatherGateway)
	weatherHandler := api.NewWeatherHandler(weatherUseCase)
	healthHandler := api.NewHealthCheck()

	webserver := web.NewWebServer(configs.WebServerPort)
	webserver.AddHandler("/", weatherHandler.GetWeather)
	webserver.AddHandler("/health", healthHandler.HealthCheck)
	webserver.Start()
}
