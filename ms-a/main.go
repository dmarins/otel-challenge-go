package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	_, span := tracer.Start(r.Context(), "handler:get:cep")

	var cepRequest CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&cepRequest); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !validateCEP(cepRequest.CEP) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		return
	}

	span.End()

	ctx, span := tracer.Start(r.Context(), "handler:call:ms-b")
	defer span.End()

	serviceBURL := "http://ms-b:8081/weather"
	reqBody, _ := json.Marshal(cepRequest)
	req, err := http.NewRequestWithContext(ctx, "POST", serviceBURL, bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

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

	http.HandleFunc("/cep", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
