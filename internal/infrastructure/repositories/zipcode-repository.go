package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmarins/otel-challenge-go/configs"
	"github.com/dmarins/otel-challenge-go/internal/domain/models"
)

type ZipcodeRepositoryInterface interface {
	GetZipcodeInfo(zipcode string) (*models.Zipcode, error)
}

type ZipcodeRepository struct {
	BaseURL    string
	HTTPClient http.Client
}

func NewZipcodeRepository(httpClient *http.Client) *ZipcodeRepository {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	return &ZipcodeRepository{
		BaseURL:    configs.ViaCepApiUrl,
		HTTPClient: *httpClient,
	}
}

func (r *ZipcodeRepository) GetZipcodeInfo(zipcode string) (*models.Zipcode, error) {
	url := fmt.Sprintf("%s/ws/%s/json", r.BaseURL, zipcode)

	resp, err := r.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result models.Zipcode
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
