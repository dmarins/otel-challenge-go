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

var setZipcodeRepository = wire.NewSet(
	repositories.NewZipcodeRepository,
	wire.Bind(new(repositories.ZipcodeRepositoryInterface), new(*repositories.ZipcodeRepository)),
)

var setWeatherRepository = wire.NewSet(
	repositories.NewWeatherRepository,
	wire.Bind(new(repositories.WeatherRepositoryInterface), new(*repositories.WeatherRepository)),
)

func NewGetWeatherByZipcodeUseCase(httpClient *http.Client) *usecase.GetWeatherByZipcodeUseCase {
	wire.Build(
		setZipcodeRepository,
		setWeatherRepository,
		usecase.NewGetWeatherByZipcodeUseCase,
	)

	return &usecase.GetWeatherByZipcodeUseCase{}
}

func NewWeatherHttpHandler(getWeatherByZipcodeUseCase usecase.GetWeatherByZipcodeUseCase) *handlers.WeatherHttpHandler {
	wire.Build(
		handlers.NewWeatherHttpHandler,
	)

	return &handlers.WeatherHttpHandler{}
}
