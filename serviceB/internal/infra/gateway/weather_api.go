package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
)

type WeatherAPI struct {
	APIKey string
}

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

type WeatherAPIResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Location struct {
	Name string `json:"name"`
}

type Current struct {
	Temp_c float64 `json:"temp_c"`
	Temp_f float64 `json:"temp_f"`
	Temp_k float64 `json:"temp_k"`
}

func NewWeatherAPI(apikey string) *WeatherAPI {
	return &WeatherAPI{APIKey: apikey}
}

func (w *WeatherAPI) getLocation(cep string) (*ViaCEPResponse, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", url.QueryEscape(cep))
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get location data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var location ViaCEPResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (w *WeatherAPI) GetCurrentWeather(cep string) (*entity.Weather, error) {
	location, err := w.getLocation(cep)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", w.APIKey, url.QueryEscape(location.Localidade))
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get location data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var weatherResponse WeatherAPIResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return nil, err
	}

	weatherData := entity.NewWeather(
		weatherResponse.Location.Name,
		weatherResponse.Current.Temp_c,
	)

	return weatherData, nil
}
