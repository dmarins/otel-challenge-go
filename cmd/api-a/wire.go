//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/dmarins/otel-challenge-go/internal/infrastructure/repositories"
	"github.com/dmarins/otel-challenge-go/internal/infrastructure/web/handlers"
	"github.com/dmarins/otel-challenge-go/internal/usecase"
	"github.com/google/wire"
)

var setApiBRepository = wire.NewSet(
	repositories.NewApiBRepository,
	wire.Bind(new(repositories.ApiBRepositoryInterface), new(*repositories.ApiBRepository)),
)

func NewGetWeatherByZipcodeInApiBUseCase(httpClient *http.Client) *usecase.GetWeatherByZipcodeInApiBUseCase {
	wire.Build(
		setApiBRepository,
		usecase.NewGetWeatherByZipcodeInApiBUseCase,
	)

	return &usecase.GetWeatherByZipcodeInApiBUseCase{}
}

func NewZipcodeHttpHandler(getWeatherByZipcodeInApiBUseCase usecase.GetWeatherByZipcodeInApiBUseCase) *handlers.ZipcodeHttpHandler {
	wire.Build(
		handlers.NewZipcodeHttpHandler,
	)

	return &handlers.ZipcodeHttpHandler{}
}
