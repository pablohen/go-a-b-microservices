package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go-a-b-microservices/pkg/apperror"
	"go-a-b-microservices/pkg/logger"
	"go-a-b-microservices/pkg/zipcode"
	"go-a-b-microservices/service-b/internal/usecase"

	"go.opentelemetry.io/otel"
)

type Handler struct {
	zipCodeUseCase *usecase.ZipCodeUseCase
	logger         logger.Logger
}

func NewHandler(zipCodeUseCase *usecase.ZipCodeUseCase, logger logger.Logger) *Handler {
	return &Handler{
		zipCodeUseCase: zipCodeUseCase,
		logger:         logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/weather", h.ProcessZipCode)
}

func (h *Handler) ProcessZipCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(r.Context(), "http.ProcessZipCode")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read request body: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}
	defer r.Body.Close()

	var request zipcode.ZipCodeRequest
	if err := json.Unmarshal(body, &request); err != nil {
		h.logger.Error("Failed to parse JSON: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Error("Invalid ZIP code: %v", err)
		writeJSONResponse(w, http.StatusUnprocessableEntity, map[string]string{"message": apperror.ErrZipCodeInvalid.Error()})
		return
	}

	response, err := h.zipCodeUseCase.ProcessZipCode(ctx, &request)
	if err != nil {
		switch err.Error() {
		case apperror.ErrZipCodeInvalid.Error():
			writeJSONResponse(w, http.StatusUnprocessableEntity, map[string]string{"message": apperror.ErrZipCodeInvalid.Error()})
		case apperror.ErrZipCodeNotFound.Error():
			writeJSONResponse(w, http.StatusNotFound, map[string]string{"message": apperror.ErrZipCodeNotFound.Error()})
		default:
			h.logger.Error("Failed to process ZIP code: %v", err)
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "internal server error"})
		}
		return
	}

	writeJSONResponse(w, http.StatusOK, response)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
