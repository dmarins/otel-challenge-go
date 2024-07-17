package usecase_test

import (
	"testing"

	mock "github.com/dmarins/otel-challenge-go/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

type TestVars struct {
	Ctrl *gomock.Controller

	ZipcodeRepository *mock.MockZipcodeRepositoryInterface
	WeatherRepository *mock.MockWeatherRepositoryInterface
}

func BuildTestVars(t *testing.T) TestVars {
	testVars := TestVars{}

	ctrl := gomock.NewController(t)
	testVars.Ctrl = ctrl

	buildRepositoryMocks(&testVars, ctrl)

	return testVars
}

func buildRepositoryMocks(testVars *TestVars, ctrl *gomock.Controller) {
	testVars.ZipcodeRepository = mock.NewMockZipcodeRepositoryInterface(ctrl)
	testVars.WeatherRepository = mock.NewMockWeatherRepositoryInterface(ctrl)
}
