package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go-a-b-microservices/pkg/apperror"
	"go-a-b-microservices/pkg/config"
	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/zipcode"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type ServiceBClientInterface interface {
	GetWeatherByZipCode(ctx context.Context, zipCode string) (*WeatherResponse, error)
}

type ServiceBClient struct {
	client  *http.Client
	cfg     *config.Config
	logger  logger.Logger
	baseURL string
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewServiceBClient(cfg *config.Config, log logger.Logger) *ServiceBClient {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	return &ServiceBClient{
		client:  httpClient,
		cfg:     cfg,
		logger:  log,
		baseURL: cfg.ServiceBURL,
	}
}

func (c *ServiceBClient) GetWeatherByZipCode(ctx context.Context, zipCode string) (*WeatherResponse, error) {
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(ctx, "repository.GetWeatherByZipCode")
	defer span.End()

	requestBody := zipcode.ZipCodeRequest{CEP: zipCode}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		c.logger.Error("Failed to marshal request: %v", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/weather", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.logger.Error("Failed to create request: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Failed to make request to Service B: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			c.logger.Error("Failed to unmarshal error response: %v", err)
			return nil, fmt.Errorf("service B returned status %d: %s", resp.StatusCode, string(body))
		}

		switch resp.StatusCode {
		case http.StatusUnprocessableEntity:
			return nil, apperror.ErrZipCodeInvalid
		case http.StatusNotFound:
			return nil, apperror.ErrZipCodeNotFound
		default:
			return nil, errors.New(errResp.Message)
		}
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		c.logger.Error("Failed to unmarshal response: %v", err)
		return nil, err
	}

	return &weatherResp, nil
}
