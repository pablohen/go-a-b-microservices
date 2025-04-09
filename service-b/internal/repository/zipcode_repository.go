package repository

import (
	"context"

	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/zipcode"
	"go-a-b-microservices/service-b/internal/adapter/clients"

	"go.opentelemetry.io/otel"
)

type ZipCodeRepository struct {
	viaCEPClient     *clients.ViaCEPClient
	weatherAPIClient *clients.WeatherAPIClient
	logger           logger.Logger
}

func NewZipCodeRepository(
	viaCEPClient *clients.ViaCEPClient,
	weatherAPIClient *clients.WeatherAPIClient,
	log logger.Logger,
) *ZipCodeRepository {
	return &ZipCodeRepository{
		viaCEPClient:     viaCEPClient,
		weatherAPIClient: weatherAPIClient,
		logger:           log,
	}
}

func (r *ZipCodeRepository) GetLocationByZipCode(ctx context.Context, zipCode string) (*zipcode.Location, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "repository.GetLocationByZipCode")
	defer span.End()

	return r.viaCEPClient.GetLocationByZipCode(ctx, zipCode)
}

func (r *ZipCodeRepository) GetWeatherByCity(ctx context.Context, city string) (*zipcode.WeatherData, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "repository.GetWeatherByCity")
	defer span.End()

	return r.weatherAPIClient.GetWeatherByCity(ctx, city)
}
