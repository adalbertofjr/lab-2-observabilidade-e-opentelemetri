package entity

type Weather struct {
	City   string
	Temp_c float64
	Temp_f float64
	Temp_k float64
}

func NewWeather(city string, tempC float64) *Weather {
	weather := Weather{
		City:   city,
		Temp_c: tempC,
	}
	weather.calcFahrenheit()
	weather.calcKelvin()

	return &weather
}

func (w *Weather) calcFahrenheit() (*Weather, error) {
	w.Temp_f = (w.Temp_c * 1.8) + 32
	return w, nil
}

func (w *Weather) calcKelvin() (*Weather, error) {
	w.Temp_k = w.Temp_c + 273
	return w, nil
}
