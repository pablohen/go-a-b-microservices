# Go A-B Microservices

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌──────────┐     ┌────────────┐
│             │     │             │     │          │     │            │
│   Client    ├────►│  Service A  ├────►│ Service B├────►│  ViaCEP    │
│             │     │             │     │          │     │            │
└─────────────┘     └─────────────┘     └──────┬───┘     └────────────┘
                                               │
                                               │
                                               ▼
                                        ┌────────────┐
                                        │            │
                                        │ WeatherAPI │
                                        │            │
                                        └────────────┘
```

## Features

- REST API for weather information based on Brazilian ZIP codes
- Microservices architecture with separate components
- External API integration (ViaCEP and WeatherAPI)
- Distributed tracing with Zipkin
- Containerization with Docker and Docker Compose
- Structured logging
- Configuration through environment variables

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- WeatherAPI API key (register at [WeatherAPI](https://www.weatherapi.com/))

## Setup

1. Clone the repository:

```bash
git clone https://github.com/pablohen/go-a-b-microservices.git
cd go-a-b-microservices
```

2. Create a `.env` file with your WeatherAPI key (see `.env.example`):

```
WEATHER_API_KEY=your_weatherapi_key_here
```

3. Build and run the services using Docker Compose:

```bash
docker compose up --build
```

This will start Service A, Service B, and Zipkin for distributed tracing.

## API Usage

### Get Weather by ZIP Code

```
POST /zipcode
```

Request body:

```json
{
  "cep": "13484000"
}
```

Response:

```json
{
  "city": "Limeira",
  "temp_C": 28.3,
  "temp_F": 82.94,
  "temp_K": 301.3
}
```

Note: The ZIP code (CEP) must be 8 digits without any special characters.

## Service Details

### Service A

- Entry point for client requests
- Validates the ZIP code format
- Forwards requests to Service B
- Handles error responses
- Exposes REST API endpoints

### Service B

- Processes business logic
- Fetches location data from ViaCEP API
- Retrieves weather information from WeatherAPI
- Calculates temperature in different units (Celsius, Fahrenheit, Kelvin)

## Project Structure

```
.
├── docker-compose.yaml         # Docker Compose configuration
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── pkg/                        # Shared packages
│   ├── apperror/               # Application error definitions
│   ├── config/                 # Configuration utilities
│   ├── logger/                 # Logging utilities
│   ├── otel/                   # OpenTelemetry integration
│   └── zipcode/                # ZIP code related structures
├── service-a/                  # Service A implementation
│   ├── Dockerfile              # Docker build instructions
│   ├── cmd/                    # Command-line entry point
│   └── internal/               # Internal packages
│       ├── adapter/            # External adapters
│       ├── repository/         # Data access layer
│       └── usecase/            # Business logic
└── service-b/                  # Service B implementation
    ├── Dockerfile              # Docker build instructions
    ├── cmd/                    # Command-line entry point
    └── internal/               # Internal packages
        ├── adapter/            # External adapters
        ├── repository/         # Data access layer
        └── usecase/            # Business logic
```
