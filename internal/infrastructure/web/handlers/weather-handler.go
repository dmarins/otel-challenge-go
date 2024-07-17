package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/dmarins/otel-challenge-go/internal/usecase"
	"github.com/go-chi/chi"
)

type WeatherResponse struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type WeatherHttpHandler struct {
	GetWeatherByZipcodeUseCase usecase.GetWeatherByZipcodeUseCase
}

func NewWeatherHttpHandler(getWeatherByZipcodeUseCase usecase.GetWeatherByZipcodeUseCase) *WeatherHttpHandler {
	return &WeatherHttpHandler{
		GetWeatherByZipcodeUseCase: getWeatherByZipcodeUseCase,
	}
}

func (handler *WeatherHttpHandler) GetWeatherInfo(w http.ResponseWriter, r *http.Request) {

	zipcode := chi.URLParam(r, "zipcode")
	if zipcode == "" {
		http.Error(w, "Invalid zipcode.", http.StatusUnprocessableEntity)
		return
	}

	matched, err := regexp.MatchString(`^\d{8}$`, zipcode)
	if err != nil {
		http.Error(w, "Invalid zipcode.", http.StatusUnprocessableEntity)
		return
	}

	if !matched {
		http.Error(w, "Invalid zipcode.", http.StatusUnprocessableEntity)
		return
	}

	dto := usecase.InputDTO{
		Zipcode: zipcode,
	}

	output, err := handler.GetWeatherByZipcodeUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if output == nil {
		http.Error(w, "Can not find zipcode.", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	result := WeatherResponse{
		Celsius:    output.Celsius,
		Fahrenheit: output.Fahrenheit,
		Kelvin:     output.Kelvin,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
