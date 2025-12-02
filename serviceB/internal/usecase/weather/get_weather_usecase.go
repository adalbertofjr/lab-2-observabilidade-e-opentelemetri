package usecase

import (
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
	domainGateway "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/gateway"
	internalerror "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/internal_error"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/pkg/utility"
)

type WeatherUseCase struct {
	weatherGateway domainGateway.WeatherGateway
}

func NewWeatherUseCase(gateway domainGateway.WeatherGateway) *WeatherUseCase {
	return &WeatherUseCase{weatherGateway: gateway}
}

func (w *WeatherUseCase) GetCurrentWeather(cep string) (*entity.Weather, *internalerror.InternalError) {
	cepFormated, err := utility.CEPFormatter(cep)
	if err != nil {
		return nil, internalerror.CEPInvalidError()
	}

	weatherData, err := w.weatherGateway.GetCurrentWeather(cepFormated)
	if err != nil {
		return nil, internalerror.CEPNotFoundError()
	}

	currentWeather := entity.NewWeather(
		weatherData.City,
		weatherData.Temp_c)

	return currentWeather, nil
}
