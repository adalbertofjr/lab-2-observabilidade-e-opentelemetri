package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/api/dto"
	"github.com/adalbertofjr/lab-2-go-service-a-otel/internal/domain/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WeatherUseCaseInterface interface {
	GetCurrentWeather(ctx context.Context, cep string) (*entity.Weather, error)
}

type WeatherHandler struct {
	usecase WeatherUseCaseInterface
	tracer  trace.Tracer
}

func NewWeatherHandler(useCase WeatherUseCaseInterface, tracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{usecase: useCase, tracer: tracer}
}

type CEP struct {
	CEP string `json:"cep"`
}

func (c *WeatherHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, spanStart := c.tracer.Start(ctx, "POST /cep")
	defer spanStart.End()

	req := r.Body
	defer req.Close()

	body, err := io.ReadAll(req)
	if err != nil {
		panic(err)
	}

	var cep CEP
	err = json.Unmarshal(body, &cep)
	if err != nil {
		panic(err)
	}

	currentWeather, err := c.usecase.GetCurrentWeather(ctx, cep.CEP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	weatherDTO := dto.NewWeatherDTO(
		currentWeather.City,
		currentWeather.Temp_c,
		currentWeather.Temp_f,
		currentWeather.Temp_k,
	)

	weatherCurrentJSON, jsonErr := json.Marshal(weatherDTO)
	if jsonErr != nil {
		http.Error(w, "Error marshalling location data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(weatherCurrentJSON))

	// spanStart.End()
}
