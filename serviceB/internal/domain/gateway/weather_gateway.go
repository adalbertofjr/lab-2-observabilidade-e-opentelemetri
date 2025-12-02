package gateway

import (
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
)

type WeatherGateway interface {
	GetCurrentWeather(cep string) (*entity.Weather, error)
}
