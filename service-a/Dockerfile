FROM golang:1.24.2-alpine AS build

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-a ./service-a/cmd/main.go

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /service-a .

# Create .env file (will be overridden by environment variables)
RUN touch .env

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./service-a"]
