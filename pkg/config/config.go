package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceAPort   string
	ServiceBPort   string
	ServiceBURL    string
	ViaCepURL      string
	WeatherAPIURL  string
	WeatherAPIKey  string
	ZipkinEndpoint string
	ServiceName    string
}

func LoadConfig(serviceName string) (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		ServiceAPort:   getEnv("SERVICE_A_PORT", "8080"),
		ServiceBPort:   getEnv("SERVICE_B_PORT", "8081"),
		ServiceBURL:    getEnv("SERVICE_B_URL", "http://localhost:8081"),
		ViaCepURL:      getEnv("VIA_CEP_URL", "https://viacep.com.br/ws"),
		WeatherAPIURL:  getEnv("WEATHER_API_URL", "https://api.weatherapi.com/v1/current.json"),
		WeatherAPIKey:  getEnv("WEATHER_API_KEY", ""),
		ZipkinEndpoint: getEnv("ZIPKIN_ENDPOINT", "http://localhost:9411/api/v2/spans"),
		ServiceName:    serviceName,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
