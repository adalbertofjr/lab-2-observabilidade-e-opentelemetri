package usecase

import (
	"errors"
	"testing"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/domain/entity"
)

// MockWeatherGateway é um mock do WeatherGateway para testes
type MockWeatherGateway struct {
	mockGetCurrentWeather func(cep string) (*entity.Weather, error)
}

func (m *MockWeatherGateway) GetCurrentWeather(cep string) (*entity.Weather, error) {
	return m.mockGetCurrentWeather(cep)
}

func TestGetCurrentWeather_Success(t *testing.T) {
	// Arrange
	mockGateway := &MockWeatherGateway{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
			return entity.NewWeather("São Paulo", 25.5), nil
		},
	}
	useCase := NewWeatherUseCase(mockGateway)

	// Act
	result, err := useCase.GetCurrentWeather("04446-160")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.City != "São Paulo" {
		t.Errorf("Expected location 'São Paulo', got '%s'", result.City)
	}
	if result.Temp_c != 25.5 {
		t.Errorf("Expected temp_c 25.5, got %.1f", result.Temp_c)
	}
	if result.Temp_f != 77.9 {
		t.Errorf("Expected temp_f 77.9, got %.1f", result.Temp_f)
	}
	if result.Temp_k != 298.5 {
		t.Errorf("Expected temp_k 298.5, got %.1f", result.Temp_k)
	}
}

func TestGetCurrentWeather_CEPWithoutDash(t *testing.T) {
	// Arrange
	mockGateway := &MockWeatherGateway{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
			if cep != "04446160" {
				t.Errorf("Expected CEP '04446160', got '%s'", cep)
			}
			return entity.NewWeather("São Paulo", 20.0), nil
		},
	}
	useCase := NewWeatherUseCase(mockGateway)

	// Act
	result, err := useCase.GetCurrentWeather("04446160")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetCurrentWeather_CEPWithDash(t *testing.T) {
	// Arrange
	mockGateway := &MockWeatherGateway{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
			if cep != "04446160" {
				t.Errorf("Expected formatted CEP '04446160', got '%s'", cep)
			}
			return entity.NewWeather("São Paulo", 20.0), nil
		},
	}
	useCase := NewWeatherUseCase(mockGateway)

	// Act
	result, err := useCase.GetCurrentWeather("04446-160")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != nil && result.City != "São Paulo" {
		t.Errorf("Expected location 'São Paulo', got '%s'", result.City)
	}
}

func TestGetCurrentWeather_InvalidCEP(t *testing.T) {
	// Arrange
	mockGateway := &MockWeatherGateway{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
			t.Error("Gateway should not be called for invalid CEP")
			return nil, errors.New("should not reach here")
		},
	}
	useCase := NewWeatherUseCase(mockGateway)

	invalidCEPs := []string{
		"1234567",
		"123456789",
		"12345-67",
		"ABCDE-FGH",
		"12a45-678",
		"",
		"999999999",
	}

	// Act & Assert
	for _, invalidCEP := range invalidCEPs {
		result, err := useCase.GetCurrentWeather(invalidCEP)

		if err == nil {
			t.Errorf("Expected error for invalid CEP '%s', got nil", invalidCEP)
			continue
		}
		if result != nil {
			t.Errorf("Expected nil result for invalid CEP '%s', got %v", invalidCEP, result)
		}
		if err.Code != 422 {
			t.Errorf("Expected error code 422 for invalid CEP '%s', got %d", invalidCEP, err.Code)
		}
		if err.MSG != "Invalid zipcode" {
			t.Errorf("Expected error message 'Invalid zipcode' for CEP '%s', got '%s'", invalidCEP, err.MSG)
		}
	}
}

func TestGetCurrentWeather_CEPNotFound(t *testing.T) {
	// Arrange
	mockGateway := &MockWeatherGateway{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
			return nil, errors.New("CEP not found in external API")
		},
	}
	useCase := NewWeatherUseCase(mockGateway)

	// Act
	result, err := useCase.GetCurrentWeather("99999-999")

	// Assert
	if err == nil {
		t.Fatal("Expected error for CEP not found, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if err.Code != 404 {
		t.Errorf("Expected error code 404, got %d", err.Code)
	}
	if err.MSG != "Can not find zipcode" {
		t.Errorf("Expected error message 'Can not find zipcode', got '%s'", err.MSG)
	}
}

func TestGetCurrentWeather_GatewayError(t *testing.T) {
	// Arrange
	mockGateway := &MockWeatherGateway{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
			return nil, errors.New("timeout connecting to weather API")
		},
	}
	useCase := NewWeatherUseCase(mockGateway)

	// Act
	result, err := useCase.GetCurrentWeather("04446-160")

	// Assert
	if err == nil {
		t.Fatal("Expected error from gateway, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if err.Code != 404 {
		t.Errorf("Expected error code 404, got %d", err.Code)
	}
	if err.MSG != "Can not find zipcode" {
		t.Errorf("Expected error message 'Can not find zipcode', got '%s'", err.MSG)
	}
}

func TestGetCurrentWeather_VariousValidCEPs(t *testing.T) {
	validCEPs := map[string]string{
		"04446-160": "04446160",
		"04446160":  "04446160",
		"00000-000": "00000000",
		"99999-999": "99999999",
		"12345678":  "12345678",
	}

	for input, expectedFormatted := range validCEPs {
		mockGateway := &MockWeatherGateway{
			mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
				if cep != expectedFormatted {
					t.Errorf("Expected formatted CEP '%s', got '%s'", expectedFormatted, cep)
				}
				return entity.NewWeather("Test City", 15.0), nil
			},
		}
		useCase := NewWeatherUseCase(mockGateway)

		result, err := useCase.GetCurrentWeather(input)

		if err != nil {
			t.Errorf("CEP '%s': Expected no error, got %v", input, err)
		}
		if result == nil {
			t.Errorf("CEP '%s': Expected result, got nil", input)
		}
	}
}

func TestGetCurrentWeather_NegativeTemperature(t *testing.T) {
	// Arrange
	mockGateway := &MockWeatherGateway{
		mockGetCurrentWeather: func(cep string) (*entity.Weather, error) {
			return entity.NewWeather("Polo Norte", -40.0), nil
		},
	}
	useCase := NewWeatherUseCase(mockGateway)

	// Act
	result, err := useCase.GetCurrentWeather("00000-000")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if result.Temp_c != -40.0 {
		t.Errorf("Expected temp_c -40.0, got %.1f", result.Temp_c)
	}
	if result.Temp_f != -40.0 {
		t.Errorf("Expected temp_f -40.0, got %.1f", result.Temp_f)
	}
	if result.Temp_k != 233.0 {
		t.Errorf("Expected temp_k 233.0, got %.1f", result.Temp_k)
	}
}
