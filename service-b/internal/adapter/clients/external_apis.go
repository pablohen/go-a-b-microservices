package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go-a-b-microservices/pkg/apperror"
	"go-a-b-microservices/pkg/config"
	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/zipcode"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type ViaCepErrorResponse struct {
	Erro string `json:"erro"`
}

type ViaCEPClient struct {
	client  *http.Client
	baseURL string
	logger  logger.Logger
}

func NewViaCEPClient(cfg *config.Config, log logger.Logger) *ViaCEPClient {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	return &ViaCEPClient{
		client:  httpClient,
		baseURL: cfg.ViaCepURL,
		logger:  log,
	}
}

func (c *ViaCEPClient) GetLocationByZipCode(ctx context.Context, zipCode string) (*zipcode.Location, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "client.ViaCEP.GetLocationByZipCode")
	defer span.End()

	url := fmt.Sprintf("%s/%s/json", c.baseURL, zipCode)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("Failed to create request: %v", err)
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request to ViaCEP: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("ViaCEP returned non-OK status: %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to get location: status %d", resp.StatusCode)
	}

	var errorResponse ViaCepErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err == nil && errorResponse.Erro == "true" {
		c.logger.Error("ViaCEP returned error for zipcode %s", zipCode)
		return nil, apperror.ErrZipCodeNotFound
	}

	var location zipcode.Location
	if err := json.Unmarshal(body, &location); err != nil {
		c.logger.Error("Failed to unmarshal response: %v", err)
		return nil, err
	}

	return &location, nil
}

type WeatherAPIClient struct {
	client  *http.Client
	baseURL string
	apiKey  string
	logger  logger.Logger
}

func NewWeatherAPIClient(cfg *config.Config, log logger.Logger) *WeatherAPIClient {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	return &WeatherAPIClient{
		client:  httpClient,
		baseURL: cfg.WeatherAPIURL,
		apiKey:  cfg.WeatherAPIKey,
		logger:  log,
	}
}

func (c *WeatherAPIClient) GetWeatherByCity(ctx context.Context, city string) (*zipcode.WeatherData, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "client.WeatherAPI.GetWeatherByCity")
	defer span.End()

	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		c.logger.Error("Failed to parse URL: %v", err)
		return nil, err
	}

	q := reqURL.Query()
	q.Set("key", c.apiKey)
	q.Set("q", city)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		c.logger.Error("Failed to create request: %v", err)
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request to WeatherAPI: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("WeatherAPI returned non-OK status: %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to get weather: status %d", resp.StatusCode)
	}

	var weather zipcode.WeatherData
	if err := json.Unmarshal(body, &weather); err != nil {
		c.logger.Error("Failed to unmarshal response: %v", err)
		return nil, err
	}

	return &weather, nil
}
