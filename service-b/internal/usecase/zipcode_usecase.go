package usecase

import (
	"context"

	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/zipcode"
	"go-a-b-microservices/service-b/internal/repository"

	"go.opentelemetry.io/otel"
)

type ZipCodeUseCase struct {
	repository repository.ZipCodeRepositoryInterface
	logger     logger.Logger
}

func NewZipCodeUseCase(repository repository.ZipCodeRepositoryInterface, logger logger.Logger) *ZipCodeUseCase {
	return &ZipCodeUseCase{
		repository: repository,
		logger:     logger,
	}
}

func (uc *ZipCodeUseCase) ProcessZipCode(ctx context.Context, request *zipcode.ZipCodeRequest) (*zipcode.WeatherResponse, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "usecase.ProcessZipCode")
	defer span.End()

	if err := request.Validate(); err != nil {
		uc.logger.Error("Invalid ZIP code: %v", err)
		return nil, err
	}

	location, err := uc.repository.GetLocationByZipCode(ctx, request.CEP)
	if err != nil {
		uc.logger.Error("Error getting location: %v", err)
		return nil, err
	}

	weather, err := uc.repository.GetWeatherByCity(ctx, location.City)
	if err != nil {
		uc.logger.Error("Error getting weather: %v", err)
		return nil, err
	}

	tempC := weather.Current.TempC
	tempF := celsiusToFahrenheit(tempC)
	tempK := celsiusToKelvin(tempC)

	response := &zipcode.WeatherResponse{
		City:  location.City,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	return response, nil
}

func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}
