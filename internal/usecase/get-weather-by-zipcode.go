package usecase

import "github.com/dmarins/otel-challenge-go/internal/infrastructure/repositories"

type (
	InputDTO struct {
		Zipcode string `json:"zipcode"`
	}

	OutputDTO struct {
		Celsius    float64
		Fahrenheit float64
		Kelvin     float64
	}

	GetWeatherByZipcodeUseCase struct {
		ZipcodeRepository repositories.ZipcodeRepositoryInterface
		WeatherRepository repositories.WeatherRepositoryInterface
	}
)

func NewGetWeatherByZipcodeUseCase(zipcodeRepository repositories.ZipcodeRepositoryInterface, weatherRepository repositories.WeatherRepositoryInterface) *GetWeatherByZipcodeUseCase {
	return &GetWeatherByZipcodeUseCase{
		ZipcodeRepository: zipcodeRepository,
		WeatherRepository: weatherRepository,
	}
}

func (usecase *GetWeatherByZipcodeUseCase) Execute(input InputDTO) (*OutputDTO, error) {

	zipcode, err := usecase.ZipcodeRepository.GetZipcodeInfo(input.Zipcode)
	if err != nil {
		return nil, err
	}
	if zipcode == nil || zipcode.Location == "" {
		return nil, nil
	}

	weather, err := usecase.WeatherRepository.GetWeatherInfo(zipcode.Location)
	if err != nil {
		return nil, err
	}
	if weather == nil {
		return nil, nil
	}

	return &OutputDTO{
		Celsius:    weather.Current.Celsius,
		Fahrenheit: weather.Current.Fahrenheit,
		Kelvin:     weather.Current.Celsius + 273,
	}, nil
}
