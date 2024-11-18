package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ip2country/internal/config"
	sut "ip2country/internal/middleware"
)

func TestRateLimitMiddleware(t *testing.T) {
	cfg := &config.Config{
		RateLimit:  1,
		BurstLimit: 1,
	}

	handler := sut.RateLimitMiddleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test rate limiting
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusTooManyRequests {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
	}
}

func TestErrorHandler(t *testing.T) {
	handler := sut.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	var resp sut.ErrorResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Error != "Internal Server Error" {
		t.Errorf("handler returned unexpected body: got %v want %v", resp.Error, "Internal Server Error")
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	sut.WriteError(rr, http.StatusBadRequest, "Bad Request")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("WriteError returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var resp sut.ErrorResponse
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp.Error != "Bad Request" {
		t.Errorf("WriteError returned unexpected body: got %v want %v", resp.Error, "Bad Request")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := sut.LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
