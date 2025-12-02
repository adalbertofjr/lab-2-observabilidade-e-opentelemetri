package entity

import "testing"

func TestNewWeather_PositiveTemperature(t *testing.T) {
	weather := NewWeather("São Paulo", 25.0)

	if weather.City != "São Paulo" {
		t.Errorf("Expected city 'São Paulo', got '%s'", weather.City)
	}
	if weather.Temp_c != 25.0 {
		t.Errorf("Expected temp_c 25.0, got %.1f", weather.Temp_c)
	}
	if weather.Temp_f != 77.0 {
		t.Errorf("Expected temp_f 77.0, got %.1f", weather.Temp_f)
	}
	if weather.Temp_k != 298.0 {
		t.Errorf("Expected temp_k 298.0, got %.1f", weather.Temp_k)
	}
}

func TestNewWeather_ZeroCelsius(t *testing.T) {
	weather := NewWeather("Polo Sul", 0.0)

	if weather.Temp_c != 0.0 {
		t.Errorf("Expected temp_c 0.0, got %.1f", weather.Temp_c)
	}
	if weather.Temp_f != 32.0 {
		t.Errorf("Expected temp_f 32.0, got %.1f", weather.Temp_f)
	}
	if weather.Temp_k != 273.0 {
		t.Errorf("Expected temp_k 273.0, got %.1f", weather.Temp_k)
	}
}

func TestNewWeather_BoilingPoint(t *testing.T) {
	weather := NewWeather("Deserto", 100.0)

	if weather.Temp_c != 100.0 {
		t.Errorf("Expected temp_c 100.0, got %.1f", weather.Temp_c)
	}
	if weather.Temp_f != 212.0 {
		t.Errorf("Expected temp_f 212.0, got %.1f", weather.Temp_f)
	}
	if weather.Temp_k != 373.0 {
		t.Errorf("Expected temp_k 373.0, got %.1f", weather.Temp_k)
	}
}

func TestNewWeather_NegativeTemperature(t *testing.T) {
	weather := NewWeather("Antártida", -40.0)

	if weather.Temp_c != -40.0 {
		t.Errorf("Expected temp_c -40.0, got %.1f", weather.Temp_c)
	}
	if weather.Temp_f != -40.0 {
		t.Errorf("Expected temp_f -40.0, got %.1f", weather.Temp_f)
	}
	if weather.Temp_k != 233.0 {
		t.Errorf("Expected temp_k 233.0, got %.1f", weather.Temp_k)
	}
}

func TestNewWeather_AbsoluteZero(t *testing.T) {
	weather := NewWeather("Espaço", -273.0)

	if weather.Temp_c != -273.0 {
		t.Errorf("Expected temp_c -273.0, got %.1f", weather.Temp_c)
	}
	expectedFahrenheit := -459.4
	if weather.Temp_f < expectedFahrenheit-0.1 || weather.Temp_f > expectedFahrenheit+0.1 {
		t.Errorf("Expected temp_f ~%.1f, got %.1f", expectedFahrenheit, weather.Temp_f)
	}
	if weather.Temp_k != 0.0 {
		t.Errorf("Expected temp_k 0.0, got %.1f", weather.Temp_k)
	}
}

func TestNewWeather_VeryHotTemperature(t *testing.T) {
	weather := NewWeather("Vulcão", 1000.0)

	if weather.Temp_c != 1000.0 {
		t.Errorf("Expected temp_c 1000.0, got %.1f", weather.Temp_c)
	}
	if weather.Temp_f != 1832.0 {
		t.Errorf("Expected temp_f 1832.0, got %.1f", weather.Temp_f)
	}
	if weather.Temp_k != 1273.0 {
		t.Errorf("Expected temp_k 1273.0, got %.1f", weather.Temp_k)
	}
}

func TestNewWeather_DecimalTemperature(t *testing.T) {
	weather := NewWeather("Rio de Janeiro", 28.7)

	if weather.Temp_c != 28.7 {
		t.Errorf("Expected temp_c 28.7, got %.1f", weather.Temp_c)
	}
	expectedFahrenheit := 83.66
	if weather.Temp_f < expectedFahrenheit-0.1 || weather.Temp_f > expectedFahrenheit+0.1 {
		t.Errorf("Expected temp_f ~%.2f, got %.2f", expectedFahrenheit, weather.Temp_f)
	}
	expectedKelvin := 301.7
	if weather.Temp_k < expectedKelvin-0.1 || weather.Temp_k > expectedKelvin+0.1 {
		t.Errorf("Expected temp_k ~%.1f, got %.1f", expectedKelvin, weather.Temp_k)
	}
}

func TestCalcFahrenheit_Accuracy(t *testing.T) {
	testCases := []struct {
		celsius    float64
		fahrenheit float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
		{25, 77},
		{37, 98.6},
		{-273, -459.4},
	}

	for _, tc := range testCases {
		weather := &Weather{Temp_c: tc.celsius}
		weather.calcFahrenheit()

		tolerance := 0.1
		if weather.Temp_f < tc.fahrenheit-tolerance || weather.Temp_f > tc.fahrenheit+tolerance {
			t.Errorf("For %.1f°C: expected %.1f°F, got %.1f°F", tc.celsius, tc.fahrenheit, weather.Temp_f)
		}
	}
}

func TestCalcKelvin_Accuracy(t *testing.T) {
	testCases := []struct {
		celsius float64
		kelvin  float64
	}{
		{0, 273},
		{100, 373},
		{-40, 233},
		{25, 298},
		{-273, 0},
		{27, 300},
	}

	for _, tc := range testCases {
		weather := &Weather{Temp_c: tc.celsius}
		weather.calcKelvin()

		if weather.Temp_k != tc.kelvin {
			t.Errorf("For %.1f°C: expected %.1fK, got %.1fK", tc.celsius, tc.kelvin, weather.Temp_k)
		}
	}
}

func TestNewWeather_LocationPreserved(t *testing.T) {
	locations := []string{
		"São Paulo",
		"Rio de Janeiro",
		"New York",
		"Tokyo",
		"",
		"Cidade com Espaços",
		"São José dos Campos",
	}

	for _, city := range locations {
		weather := NewWeather(city, 20.0)
		if weather.City != city {
			t.Errorf("Expected city '%s', got '%s'", city, weather.City)
		}
	}
}

func TestNewWeather_ReturnsNonNilPointer(t *testing.T) {
	weather := NewWeather("Test", 15.0)

	if weather == nil {
		t.Error("Expected non-nil pointer, got nil")
	}
}
