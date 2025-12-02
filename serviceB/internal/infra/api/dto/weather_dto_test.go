package dto

import (
	"encoding/json"
	"testing"
)

func TestNewWeatherDTO(t *testing.T) {
	dto := NewWeatherDTO("São Paulo", 25.0, 77.0, 298.0)

	if dto == nil {
		t.Fatal("Expected non-nil DTO, got nil")
	}
	if dto.City != "São Paulo" {
		t.Errorf("Expected city 'São Paulo', got '%s'", dto.City)
	}
	if dto.Temp_c != 25.0 {
		t.Errorf("Expected temp_c 25.0, got %.1f", dto.Temp_c)
	}
	if dto.Temp_f != 77.0 {
		t.Errorf("Expected temp_f 77.0, got %.1f", dto.Temp_f)
	}
	if dto.Temp_k != 298.0 {
		t.Errorf("Expected temp_k 298.0, got %.1f", dto.Temp_k)
	}
}

func TestWeatherDTO_JSONMarshaling(t *testing.T) {
	dto := NewWeatherDTO("Rio de Janeiro", 30.5, 86.9, 303.5)

	jsonData, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Failed to marshal DTO: %v", err)
	}

	expectedJSON := `{"city":"Rio de Janeiro","temp_c":30.5,"temp_f":86.9,"temp_k":303.5}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, string(jsonData))
	}
}

func TestWeatherDTO_JSONUnmarshaling(t *testing.T) {
	jsonData := `{"city":"Curitiba","temp_c":15.0,"temp_f":59.0,"temp_k":288.0}`

	var dto WeatherDTO
	err := json.Unmarshal([]byte(jsonData), &dto)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if dto.City != "Curitiba" {
		t.Errorf("Expected city 'Curitiba', got '%s'", dto.City)
	}
	if dto.Temp_c != 15.0 {
		t.Errorf("Expected temp_c 15.0, got %.1f", dto.Temp_c)
	}
	if dto.Temp_f != 59.0 {
		t.Errorf("Expected temp_f 59.0, got %.1f", dto.Temp_f)
	}
	if dto.Temp_k != 288.0 {
		t.Errorf("Expected temp_k 288.0, got %.1f", dto.Temp_k)
	}
}

func TestWeatherDTO_JSONFieldNames(t *testing.T) {
	dto := NewWeatherDTO("Test", 20.0, 68.0, 293.0)

	jsonData, _ := json.Marshal(dto)
	var rawJSON map[string]interface{}
	json.Unmarshal(jsonData, &rawJSON)

	expectedFields := []string{"city", "temp_c", "temp_f", "temp_k"}
	for _, field := range expectedFields {
		if _, exists := rawJSON[field]; !exists {
			t.Errorf("Expected JSON field '%s' not found", field)
		}
	}

	if len(rawJSON) != 4 {
		t.Errorf("Expected 4 JSON fields, got %d", len(rawJSON))
	}
}

func TestWeatherDTO_DifferentValues(t *testing.T) {
	testCases := []struct {
		name     string
		location string
		tempC    float64
		tempF    float64
		tempK    float64
	}{
		{"Zero Celsius", "Polo Sul", 0.0, 32.0, 273.0},
		{"Negative", "Sibéria", -40.0, -40.0, 233.0},
		{"Hot", "Deserto", 50.0, 122.0, 323.0},
		{"Decimal", "São Paulo", 22.5, 72.5, 295.5},
		{"Empty Location", "", 20.0, 68.0, 293.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dto := NewWeatherDTO(tc.location, tc.tempC, tc.tempF, tc.tempK)

			if dto.City != tc.location {
				t.Errorf("Expected location '%s', got '%s'", tc.location, dto.City)
			}
			if dto.Temp_c != tc.tempC {
				t.Errorf("Expected temp_c %.1f, got %.1f", tc.tempC, dto.Temp_c)
			}
			if dto.Temp_f != tc.tempF {
				t.Errorf("Expected temp_f %.1f, got %.1f", tc.tempF, dto.Temp_f)
			}
			if dto.Temp_k != tc.tempK {
				t.Errorf("Expected temp_k %.1f, got %.1f", tc.tempK, dto.Temp_k)
			}
		})
	}
}

func TestWeatherDTO_JSONRoundTrip(t *testing.T) {
	original := NewWeatherDTO("Brasília", 28.0, 82.4, 301.0)

	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled WeatherDTO
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshaled.City != original.City {
		t.Errorf("City changed after round trip")
	}
	if unmarshaled.Temp_c != original.Temp_c {
		t.Errorf("Temp_c changed after round trip")
	}
	if unmarshaled.Temp_f != original.Temp_f {
		t.Errorf("Temp_f changed after round trip")
	}
	if unmarshaled.Temp_k != original.Temp_k {
		t.Errorf("Temp_k changed after round trip")
	}
}
