package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var tracer trace.Tracer

func fetchLocation(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get location")
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["erro"] != nil {
		return "", fmt.Errorf("can not find zipcode")
	}
	return result["localidade"].(string), nil
}

func fetchWeather(ctx context.Context, city string) (float64, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=a9526be837464f3b82814230241307&q=%s&aqi=no", url.QueryEscape(city))
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get weather")
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["current"].(map[string]interface{})["temp_c"].(float64), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "handler::receives::cep")

	var cepRequest CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&cepRequest); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	span.End()

	ctx, span := tracer.Start(r.Context(), "handler::call::cep-api")

	city, err := fetchLocation(ctx, cepRequest.CEP)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "can not find zipcode"})
		return
	}

	span.End()

	ctx, span = tracer.Start(r.Context(), "handler::call::weather-api")

	tempC, err := fetchWeather(ctx, city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "failed to get weather"})
		return
	}

	span.End()

	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	response := WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	tracer = tp.Tracer("github.com/dmarins/otel-challenge-go")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/weather", handler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
