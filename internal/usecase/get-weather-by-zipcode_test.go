package usecase_test

import (
	"errors"
	"testing"

	"github.com/dmarins/otel-challenge-go/internal/domain/models"
	"github.com/dmarins/otel-challenge-go/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func buildTestContext(t *testing.T) (usecase.GetWeatherByZipcodeUseCase, TestVars) {
	testVars := BuildTestVars(t)
	sut := usecase.NewGetWeatherByZipcodeUseCase(testVars.ZipcodeRepository, testVars.WeatherRepository)

	return *sut, testVars
}

func TestGetWeatherByZipcodeUseCaseExecute_WhenZipcodeRepositoryFails(t *testing.T) {

	sut, testVars := buildTestContext(t)

	testVars.
		ZipcodeRepository.
		EXPECT().
		GetZipcodeInfo(gomock.Any()).
		Return(nil, errors.New("fail"))

	result, err := sut.Execute(usecase.InputDTO{Zipcode: "12345678"})

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestGetWeatherByZipcodeUseCaseExecute_WhenZipcodeRepositoryReturnsNil(t *testing.T) {

	sut, testVars := buildTestContext(t)

	testVars.
		ZipcodeRepository.
		EXPECT().
		GetZipcodeInfo(gomock.Any()).
		Return(nil, nil)

	result, err := sut.Execute(usecase.InputDTO{Zipcode: "12345678"})

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestGetWeatherByZipcodeUseCaseExecute_WhenZipcodeRepositoryReturnsEmpty(t *testing.T) {

	sut, testVars := buildTestContext(t)

	testVars.
		ZipcodeRepository.
		EXPECT().
		GetZipcodeInfo(gomock.Any()).
		Return(&models.Zipcode{Location: ""}, nil)

	result, err := sut.Execute(usecase.InputDTO{Zipcode: "12345678"})

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestGetWeatherByZipcodeUseCaseExecute_WhenWeatherRepositoryFails(t *testing.T) {

	sut, testVars := buildTestContext(t)

	testVars.
		ZipcodeRepository.
		EXPECT().
		GetZipcodeInfo(gomock.Any()).
		Return(&models.Zipcode{Location: "Londres"}, nil)

	testVars.
		WeatherRepository.
		EXPECT().
		GetWeatherInfo(gomock.Any()).
		Return(nil, errors.New("fail"))

	result, err := sut.Execute(usecase.InputDTO{Zipcode: "12345678"})

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestGetWeatherByWeatherUseCaseExecute_WhenWeatherRepositoryReturnsNil(t *testing.T) {

	sut, testVars := buildTestContext(t)

	testVars.
		ZipcodeRepository.
		EXPECT().
		GetZipcodeInfo(gomock.Any()).
		Return(&models.Zipcode{Location: "Londres"}, nil)

	testVars.
		WeatherRepository.
		EXPECT().
		GetWeatherInfo(gomock.Any()).
		Return(nil, nil)

	result, err := sut.Execute(usecase.InputDTO{Zipcode: "12345678"})

	assert.Nil(t, result)
	assert.Nil(t, err)
}

func TestGetWeatherByWeatherUseCaseExecute_ShouldBeAsExpected(t *testing.T) {

	sut, testVars := buildTestContext(t)

	testVars.
		ZipcodeRepository.
		EXPECT().
		GetZipcodeInfo(gomock.Any()).
		Return(&models.Zipcode{Location: "Londres"}, nil)

	testVars.
		WeatherRepository.
		EXPECT().
		GetWeatherInfo(gomock.Any()).
		Return(&models.Weather{Current: models.Current{Celsius: 22.1, Fahrenheit: 71.78}}, nil)

	result, err := sut.Execute(usecase.InputDTO{Zipcode: "12345678"})

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Celsius, 22.1)
	assert.Equal(t, result.Fahrenheit, result.Celsius*1.8+32.0)
	assert.Equal(t, result.Kelvin, result.Celsius+273.0)
}
