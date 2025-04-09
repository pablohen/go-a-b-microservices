package usecase

import (
	"context"

	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/zipcode"
	"go-a-b-microservices/service-a/internal/repository"

	"go.opentelemetry.io/otel"
)

type ZipCodeUseCase struct {
	serviceBClient *repository.ServiceBClient
	logger         logger.Logger
}

func NewZipCodeUseCase(serviceBClient *repository.ServiceBClient, logger logger.Logger) *ZipCodeUseCase {
	return &ZipCodeUseCase{
		serviceBClient: serviceBClient,
		logger:         logger,
	}
}

func (uc *ZipCodeUseCase) ProcessZipCode(ctx context.Context, request *zipcode.ZipCodeRequest) (*repository.WeatherResponse, error) {
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(ctx, "usecase.ProcessZipCode")
	defer span.End()

	if err := request.Validate(); err != nil {
		uc.logger.Error("Invalid ZIP code: %v", err)
		return nil, err
	}

	response, err := uc.serviceBClient.GetWeatherByZipCode(ctx, request.CEP)
	if err != nil {
		uc.logger.Error("Error getting weather information: %v", err)
		return nil, err
	}

	return response, nil
}
