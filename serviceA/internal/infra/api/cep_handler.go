package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/adalbertofjr/lab-2-go-service-a-otel/pkg/utility"
)

type CEPHandler struct {
	CEP string `json:"cep"`
}

func NewCEPHandler() *CEPHandler {
	return &CEPHandler{}
}

func (c *CEPHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	req := r.Body
	defer req.Close()

	body, err := io.ReadAll(req)
	if err != nil {
		panic(err)
	}
	var cep CEPHandler
	err = json.Unmarshal(body, &cep)
	if err != nil {
		panic(err)
	}

	cepValidated, err := utility.CEPValidator(cep.CEP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	httpClient := &http.Client{}
	reqWeather, err := http.NewRequest("GET", "http://localhost:8000/?cep="+cepValidated, nil)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(reqWeather)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyWeather, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyWeather, &cep)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyWeather)
}
