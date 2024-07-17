package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/dmarins/otel-challenge-go/internal/usecase"
	"github.com/go-chi/chi"
)

type ZipcodeResponse struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type ZipcodeHttpHandler struct {
	GetWeatherByZipcodeInApiBUseCase usecase.GetWeatherByZipcodeInApiBUseCase
}

func NewZipcodeHttpHandler(getWeatherByZipcodeInApiBUseCase usecase.GetWeatherByZipcodeInApiBUseCase) *ZipcodeHttpHandler {
	return &ZipcodeHttpHandler{
		GetWeatherByZipcodeInApiBUseCase: getWeatherByZipcodeInApiBUseCase,
	}
}

func (handler *ZipcodeHttpHandler) PostZipcode(w http.ResponseWriter, r *http.Request) {

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

	dto := usecase.InputDTOApiB{
		Zipcode: zipcode,
	}

	output, err := handler.GetWeatherByZipcodeInApiBUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	result := ZipcodeResponse{
		City:       output.City,
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
