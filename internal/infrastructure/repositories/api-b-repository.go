package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmarins/otel-challenge-go/configs"
	"github.com/dmarins/otel-challenge-go/internal/domain/models"
)

type ApiBRepositoryInterface interface {
	GetZipcodeInApiB(zipcode string) (*models.ZipcodeApiB, error)
}

type ApiBRepository struct {
	BaseURL    string
	HTTPClient http.Client
}

func NewApiBRepository(httpClient *http.Client) *ApiBRepository {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	return &ApiBRepository{
		BaseURL:    configs.MicroserviceBApiUrl,
		HTTPClient: *httpClient,
	}
}

func (r *ApiBRepository) GetZipcodeInApiB(zipcode string) (*models.ZipcodeApiB, error) {
	url := fmt.Sprintf("%s/weather/%s", r.BaseURL, zipcode)

	resp, err := r.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result models.ZipcodeApiB
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
