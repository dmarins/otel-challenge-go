package main

import (
	"fmt"
	"net/http"

	"github.com/dmarins/otel-challenge-go/configs"
	"github.com/dmarins/otel-challenge-go/internal/infrastructure/web/server"
)

func main() {
	// Envs
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	httpClient := http.DefaultClient
	getWeatherByZipcodeUseCase := NewGetWeatherByZipcodeUseCase(httpClient)
	weatherHttpHandler := NewWeatherHttpHandler(*getWeatherByZipcodeUseCase)

	// Http Server
	webserver := server.NewWebServer(configs.WebServerPort)

	webserver.AddHandler("get", "/weather/{zipcode}", weatherHttpHandler.GetWeatherInfo)

	fmt.Println("Starting HTTP server on port", configs.WebServerPort)

	webserver.Start()
}
