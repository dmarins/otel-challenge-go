package usecase

import "github.com/dmarins/otel-challenge-go/internal/infrastructure/repositories"

type (
	InputDTOApiB struct {
		Zipcode string `json:"zipcode"`
	}

	OutputDTOApiB struct {
		City       string
		Celsius    float64
		Fahrenheit float64
		Kelvin     float64
	}

	GetWeatherByZipcodeInApiBUseCase struct {
		ApiBRepository repositories.ApiBRepositoryInterface
	}
)

func NewGetWeatherByZipcodeInApiBUseCase(apiBRepository repositories.ApiBRepositoryInterface) *GetWeatherByZipcodeInApiBUseCase {
	return &GetWeatherByZipcodeInApiBUseCase{
		ApiBRepository: apiBRepository,
	}
}

func (usecase *GetWeatherByZipcodeInApiBUseCase) Execute(input InputDTOApiB) (*OutputDTOApiB, error) {

	zipcode, err := usecase.ApiBRepository.GetZipcodeInApiB(input.Zipcode)
	if err != nil {
		return nil, err
	}
	if zipcode == nil {
		return nil, nil
	}

	return &OutputDTOApiB{
		City:       zipcode.City,
		Celsius:    zipcode.Celsius,
		Fahrenheit: zipcode.Fahrenheit,
		Kelvin:     zipcode.Celsius,
	}, nil
}
