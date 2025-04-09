package usecase

import (
	"context"
	"errors"
	"testing"

	"go-a-b-microservices/pkg/apperror"
	"go-a-b-microservices/pkg/zipcode"
	"go-a-b-microservices/service-a/internal/repository"
)

type MockServiceBClient struct {
	GetWeatherByZipCodeFunc func(ctx context.Context, zipCode string) (*repository.WeatherResponse, error)
}

func (m *MockServiceBClient) GetWeatherByZipCode(ctx context.Context, zipCode string) (*repository.WeatherResponse, error) {
	return m.GetWeatherByZipCodeFunc(ctx, zipCode)
}

type MockLogger struct{}

func (m *MockLogger) Info(message string, args ...interface{})  {}
func (m *MockLogger) Error(message string, args ...interface{}) {}
func (m *MockLogger) Debug(message string, args ...interface{}) {}

func TestZipCodeUseCase_ProcessZipCode(t *testing.T) {
	validZipCode := "13484000"
	validCity := "Limeira"
	tempC := 28.3
	tempF := 82.94
	tempK := 301.3

	tests := []struct {
		name          string
		zipCode       string
		mockResponse  *repository.WeatherResponse
		mockError     error
		expectedCity  string
		expectedTempC float64
		expectedTempF float64
		expectedTempK float64
		expectedErr   error
	}{
		{
			name:    "valid zip code returns weather data",
			zipCode: validZipCode,
			mockResponse: &repository.WeatherResponse{
				City:  validCity,
				TempC: tempC,
				TempF: tempF,
				TempK: tempK,
			},
			mockError:     nil,
			expectedCity:  validCity,
			expectedTempC: tempC,
			expectedTempF: tempF,
			expectedTempK: tempK,
			expectedErr:   nil,
		},
		{
			name:         "invalid zip code returns error",
			zipCode:      "invalid",
			mockResponse: nil,
			mockError:    nil,
			expectedErr:  apperror.ErrZipCodeInvalid,
		},
		{
			name:         "service B client error returns error",
			zipCode:      validZipCode,
			mockResponse: nil,
			mockError:    errors.New("service client error"),
			expectedErr:  errors.New("service client error"),
		},
		{
			name:         "zip code not found returns error",
			zipCode:      validZipCode,
			mockResponse: nil,
			mockError:    apperror.ErrZipCodeNotFound,
			expectedErr:  apperror.ErrZipCodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockServiceBClient{
				GetWeatherByZipCodeFunc: func(ctx context.Context, zipCode string) (*repository.WeatherResponse, error) {
					return tt.mockResponse, tt.mockError
				},
			}
			mockLogger := &MockLogger{}
			useCase := NewZipCodeUseCase(mockClient, mockLogger)

			request := &zipcode.ZipCodeRequest{CEP: tt.zipCode}
			ctx := context.Background()

			response, err := useCase.ProcessZipCode(ctx, request)

			if tt.name == "invalid zip code returns error" {
				if err != apperror.ErrZipCodeInvalid {
					t.Errorf("Expected ErrZipCodeInvalid, but got %v", err)
				}
				return
			} else if tt.name == "zip code not found returns error" {
				if err != apperror.ErrZipCodeNotFound {
					t.Errorf("Expected ErrZipCodeNotFound, but got %v", err)
				}
				return
			} else if tt.expectedErr != nil {
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
