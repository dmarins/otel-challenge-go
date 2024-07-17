package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dmarins/otel-challenge-go/configs"
	"github.com/dmarins/otel-challenge-go/internal/domain/models"
)

type WeatherRepositoryInterface interface {
	GetWeatherInfo(location string) (*models.Weather, error)
}

type WeatherRepository struct {
	BaseURL    string
	HTTPClient http.Client
}

func NewWeatherRepository(httpClient *http.Client) *WeatherRepository {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	return &WeatherRepository{
		BaseURL:    configs.WeatherApiUrl,
		HTTPClient: *httpClient,
	}
}

func (r *WeatherRepository) GetWeatherInfo(location string) (*models.Weather, error) {
	url := fmt.Sprintf("%s/v1/current.json?q=%s&key=a9526be837464f3b82814230241307&aqi=no", r.BaseURL, url.QueryEscape(location))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	var result models.Weather
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
