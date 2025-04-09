package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-a-b-microservices/pkg/config"
	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/otel"
	custom_http "go-a-b-microservices/service-a/internal/adapter/http"
	"go-a-b-microservices/service-a/internal/repository"
	"go-a-b-microservices/service-a/internal/usecase"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	log := logger.NewLogger()
	log.Info("Starting Service A")

	cfg, err := config.LoadConfig("service-a")
	if err != nil {
		log.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	tp, err := otel.InitTracer(cfg.ServiceName, cfg.ZipkinEndpoint, log)
	if err != nil {
		log.Error("Failed to initialize tracer: %v", err)
		os.Exit(1)
	}
	defer otel.ShutdownTracer(ctx, tp, log)

	serviceBClient := repository.NewServiceBClient(cfg, log)
	zipCodeUseCase := usecase.NewZipCodeUseCase(serviceBClient, log)
	handler := custom_http.NewHandler(zipCodeUseCase, log)

	mux := http.NewServeMux()

	handler.RegisterRoutes(mux)

	otelHandler := otelhttp.NewHandler(mux, cfg.ServiceName)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServiceAPort),
		Handler: otelHandler,
	}

	go func() {
		log.Info("Service A listening on port %s", cfg.ServiceAPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Service A")
}
