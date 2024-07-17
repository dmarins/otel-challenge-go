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
	getWeatherByZipcodeInApiBUseCase := NewGetWeatherByZipcodeInApiBUseCase(httpClient)
	zipcodeHttpHandler := NewZipcodeHttpHandler(*getWeatherByZipcodeInApiBUseCase)

	// Http Server
	webserver := server.NewWebServer(configs.WebServerPort)

	webserver.AddHandler("post", "/zipcode", zipcodeHttpHandler.PostZipcode)

	fmt.Println("Starting HTTP server on port", configs.WebServerPort)

	webserver.Start()
}
