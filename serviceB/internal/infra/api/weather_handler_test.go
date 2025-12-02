package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
	internalerror "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/internal_error"
)

type MockWeatherUseCase struct {
	mockGetCurrentWeather func(cep string) (*entity.Weather, *internalerror.InternalError)
}

func (m *MockWeatherUseCase) GetCurrentWeather(cep string) (*entity.Weather, *internalerror.InternalError) {
	return m.mockGetCurrentWeather(cep)
}

func TestGetWeather_Success(t *testing.T) {
	mockUseCase := &MockWeatherUseCase{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, *internalerror.InternalError) {
			return entity.NewWeather("São Paulo", 25.5), nil
		},
	}
	handler := NewWeatherHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=04446-160", nil)
	w := httptest.NewRecorder()

	handler.GetWeather(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["city"] != "São Paulo" {
		t.Errorf("Expected city 'São Paulo', got '%v'", response["city"])
	}
	if response["temp_c"] != 25.5 {
		t.Errorf("Expected temp_c 25.5, got %v", response["temp_c"])
	}
	if response["temp_f"] != 77.9 {
		t.Errorf("Expected temp_f 77.9, got %v", response["temp_f"])
	}
	if response["temp_k"] != 298.5 {
		t.Errorf("Expected temp_k 298.5, got %v", response["temp_k"])
	}
}

func TestGetWeather_InvalidCEP(t *testing.T) {
	mockUseCase := &MockWeatherUseCase{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, *internalerror.InternalError) {
			return nil, internalerror.CEPInvalidError()
		},
	}
	handler := NewWeatherHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=123", nil)
	w := httptest.NewRecorder()

	handler.GetWeather(w, req)

	if w.Code != 422 {
		t.Errorf("Expected status code 422, got %d", w.Code)
	}

	body := w.Body.String()
	expectedMessage := "Invalid zipcode\n"
	if body != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, body)
	}
}

func TestGetWeather_CEPNotFound(t *testing.T) {
	mockUseCase := &MockWeatherUseCase{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, *internalerror.InternalError) {
			return nil, internalerror.CEPNotFoundError()
		},
	}
	handler := NewWeatherHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=99999-999", nil)
	w := httptest.NewRecorder()

	handler.GetWeather(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status code 404, got %d", w.Code)
	}

	body := w.Body.String()
	expectedMessage := "Can not find zipcode\n"
	if body != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, body)
	}
}

func TestGetWeather_MissingCEPParameter(t *testing.T) {
	mockUseCase := &MockWeatherUseCase{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, *internalerror.InternalError) {
			if cep == "" {
				return nil, internalerror.CEPInvalidError()
			}
			return entity.NewWeather("Test", 20.0), nil
		},
	}
	handler := NewWeatherHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	w := httptest.NewRecorder()

	handler.GetWeather(w, req)

	if w.Code != 422 {
		t.Errorf("Expected status code 422, got %d", w.Code)
	}
}

func TestGetWeather_JSONStructure(t *testing.T) {
	mockUseCase := &MockWeatherUseCase{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, *internalerror.InternalError) {
			return entity.NewWeather("Rio de Janeiro", 30.0), nil
		},
	}
	handler := NewWeatherHandler(mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=20000-000", nil)
	w := httptest.NewRecorder()

	handler.GetWeather(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	requiredFields := []string{"city", "temp_c", "temp_f", "temp_k"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Response missing required field: %s", field)
		}
	}

	if len(response) != 4 {
		t.Errorf("Expected 4 fields in response, got %d", len(response))
	}
}

func TestGetWeather_DifferentTemperatures(t *testing.T) {
	testCases := []struct {
		name  string
		city  string
		tempC float64
		tempF float64
		tempK float64
	}{
		{"Zero Celsius", "Polo Sul", 0.0, 32.0, 273.0},
		{"Negative Temperature", "Sibéria", -40.0, -40.0, 233.0},
		{"Hot Temperature", "Deserto", 45.0, 113.0, 318.0},
		{"Room Temperature", "Casa", 20.0, 68.0, 293.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUseCase := &MockWeatherUseCase{
				mockGetCurrentWeather: func(cep string) (*entity.Weather, *internalerror.InternalError) {
					return entity.NewWeather(tc.city, tc.tempC), nil
				},
			}
			handler := NewWeatherHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/weather?cep=12345-678", nil)
			w := httptest.NewRecorder()

			handler.GetWeather(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if response["city"] != tc.city {
				t.Errorf("Expected city '%s', got '%v'", tc.city, response["city"])
			}
			if response["temp_c"] != tc.tempC {
				t.Errorf("Expected temp_c %.1f, got %v", tc.tempC, response["temp_c"])
			}
			if response["temp_f"] != tc.tempF {
				t.Errorf("Expected temp_f %.1f, got %v", tc.tempF, response["temp_f"])
			}
			if response["temp_k"] != tc.tempK {
				t.Errorf("Expected temp_k %.1f, got %v", tc.tempK, response["temp_k"])
			}
		})
	}
}

func TestGetWeather_CEPFormats(t *testing.T) {
	cepFormats := []string{
		"04446-160",
		"04446160",
		"12345-678",
		"00000-000",
	}

	for _, cep := range cepFormats {
		t.Run("CEP_"+cep, func(t *testing.T) {
			mockUseCase := &MockWeatherUseCase{
				mockGetCurrentWeather: func(receivedCEP string) (*entity.Weather, *internalerror.InternalError) {
					return entity.NewWeather("Test City", 22.0), nil
				},
			}
			handler := NewWeatherHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/weather?cep="+cep, nil)
			w := httptest.NewRecorder()

			handler.GetWeather(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("CEP '%s': Expected status 200, got %d", cep, w.Code)
			}
		})
	}
}
