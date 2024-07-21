package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var tracer trace.Tracer

func validateCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

func handler(w http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get("x-request-id")

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	var cepRequest CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&cepRequest); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	values := []attribute.KeyValue{
		attribute.String("RequestId", requestId),
		attribute.String("CEP", cepRequest.CEP),
	}
	_, span := tracer.Start(ctx, "handler::receives::cep", trace.WithAttributes(values...))

	if !validateCEP(cepRequest.CEP) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		return
	}

	span.End()

	ctx, span = tracer.Start(ctx, "handler::post::ms-b")
	defer span.End()

	serviceBURL := "http://ms-b:8081/weather"
	reqBody, _ := json.Marshal(cepRequest)
	req, err := http.NewRequestWithContext(ctx, "POST", serviceBURL, bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-request-id", requestId)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "failed to get response from service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	tracer = tp.Tracer("github.com/dmarins/otel-challenge-go")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/cep", handler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
