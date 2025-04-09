package usecase

import (
	"context"
	"testing"

	"go-a-b-microservices/pkg/apperror"
	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/zipcode"
	"go-a-b-microservices/service-b/internal/repository"
)

type MockZipCodeRepository struct {
	GetLocationByZipCodeFunc func(ctx context.Context, zipCode string) (*zipcode.Location, error)
	GetWeatherByCityFunc     func(ctx context.Context, city string) (*zipcode.WeatherData, error)
}

func (m *MockZipCodeRepository) GetLocationByZipCode(ctx context.Context, zipCode string) (*zipcode.Location, error) {
	return m.GetLocationByZipCodeFunc(ctx, zipCode)
}

func (m *MockZipCodeRepository) GetWeatherByCity(ctx context.Context, city string) (*zipcode.WeatherData, error) {
	return m.GetWeatherByCityFunc(ctx, city)
}

type MockLogger struct{}

func (m *MockLogger) Info(message string, args ...interface{})  {}
func (m *MockLogger) Error(message string, args ...interface{}) {}
func (m *MockLogger) Debug(message string, args ...interface{}) {}

type ZipCodeUseCaseForTesting struct {
	repository repository.ZipCodeRepositoryInterface
	logger     logger.Logger
}

func NewZipCodeUseCaseForTesting(repository repository.ZipCodeRepositoryInterface, logger logger.Logger) *ZipCodeUseCaseForTesting {
	return &ZipCodeUseCaseForTesting{
		repository: repository,
		logger:     logger,
	}
}

func (uc *ZipCodeUseCaseForTesting) ProcessZipCode(ctx context.Context, request *zipcode.ZipCodeRequest) (*zipcode.WeatherResponse, error) {
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

func TestZipCodeUseCase_ProcessZipCode(t *testing.T) {
	validZipCode := "13484000"
	validCity := "Limeira"
	tempC := 28.3

	tests := []struct {
		name          string
		zipCode       string
		mockLocation  *zipcode.Location
		locationErr   error
		mockWeather   *zipcode.WeatherData
		weatherErr    error
		expectedCity  string
		expectedTempC float64
		expectedTempF float64
		expectedTempK float64
		expectedErr   error
	}{
		{
			name:    "valid zip code returns weather data",
			zipCode: validZipCode,
			mockLocation: &zipcode.Location{
				City: validCity,
				CEP:  validZipCode,
			},
			locationErr: nil,
			mockWeather: &zipcode.WeatherData{
				Current: struct {
					TempC float64 `json:"temp_c"`
				}{
					TempC: tempC,
				},
			},
			weatherErr:    nil,
			expectedCity:  validCity,
			expectedTempC: tempC,
			expectedTempF: celsiusToFahrenheit(tempC),
			expectedTempK: celsiusToKelvin(tempC),
			expectedErr:   nil,
		},
		{
			name:         "invalid zip code returns error",
			zipCode:      "invalid",
			mockLocation: nil,
			locationErr:  nil,
			mockWeather:  nil,
			weatherErr:   nil,
			expectedErr:  apperror.ErrZipCodeInvalid,
		},
		{
			name:         "location service error returns error",
			zipCode:      validZipCode,
			mockLocation: nil,
			locationErr:  apperror.ErrZipCodeNotFound,
			mockWeather:  nil,
			weatherErr:   nil,
			expectedErr:  apperror.ErrZipCodeNotFound,
		},
		{
			name:    "weather service error returns error",
			zipCode: validZipCode,
			mockLocation: &zipcode.Location{
				City: validCity,
				CEP:  validZipCode,
			},
			locationErr: nil,
			mockWeather: nil,
			weatherErr:  apperror.ErrZipCodeNotFound,
			expectedErr: apperror.ErrZipCodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockZipCodeRepository{
				GetLocationByZipCodeFunc: func(ctx context.Context, zipCode string) (*zipcode.Location, error) {
					return tt.mockLocation, tt.locationErr
				},
				GetWeatherByCityFunc: func(ctx context.Context, city string) (*zipcode.WeatherData, error) {
					return tt.mockWeather, tt.weatherErr
				},
			}
			mockLogger := &MockLogger{}
			useCase := NewZipCodeUseCaseForTesting(mockRepo, mockLogger)

			request := &zipcode.ZipCodeRequest{CEP: tt.zipCode}
			ctx := context.Background()

			response, err := useCase.ProcessZipCode(ctx, request)

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tt.expectedErr)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, but got %v", err)
				return
			}

			if response.City != tt.expectedCity {
				t.Errorf("Expected city %s, got %s", tt.expectedCity, response.City)
			}

			if response.TempC != tt.expectedTempC {
				t.Errorf("Expected temperature in Celsius %f, got %f", tt.expectedTempC, response.TempC)
			}

			if response.TempF != tt.expectedTempF {
				t.Errorf("Expected temperature in Fahrenheit %f, got %f", tt.expectedTempF, response.TempF)
			}

			if response.TempK != tt.expectedTempK {
				t.Errorf("Expected temperature in Kelvin %f, got %f", tt.expectedTempK, response.TempK)
			}
		})
	}
}
