package api

import (
	"encoding/json"
	"net/http"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/api/dto"
	internalerror "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/internal_error"
)

type WeatherUseCaseInterface interface {
	GetCurrentWeather(cep string) (*entity.Weather, *internalerror.InternalError)
}

type WeatherHandler struct {
	usecase WeatherUseCaseInterface
}

func NewWeatherHandler(useCase WeatherUseCaseInterface) *WeatherHandler {
	return &WeatherHandler{usecase: useCase}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	weatherCurrent, err := h.usecase.GetCurrentWeather(cep)
	if err != nil {
		http.Error(w, err.MSG, err.Code)
		return
	}

	weatherDTO := dto.NewWeatherDTO(
		weatherCurrent.City,
		weatherCurrent.Temp_c,
		weatherCurrent.Temp_f,
		weatherCurrent.Temp_k,
	)

	weatherCurrentJSON, jsonErr := json.Marshal(weatherDTO)
	if jsonErr != nil {
		http.Error(w, "Error marshalling location data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(weatherCurrentJSON))
}
