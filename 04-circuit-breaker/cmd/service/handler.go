package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"load-testing/internal/pkg/circuitbreaker"
)

type Handler struct {
	httpClient *http.Client
	retryCount int

	breaker *circuitbreaker.CircuitBreaker
}

func FromEnv() (*Handler, error) {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT_MS"))
	if err != nil {
		return nil, err
	}

	retryCount, err := strconv.Atoi(os.Getenv("RETRY_COUNT"))
	if err != nil {
		return nil, err
	}

	cli := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	maxErrorPerSecond, err := strconv.ParseFloat(os.Getenv("BREAKER_MAX_ERRORS_PER_SECOND"), 64)
	if err != nil {
		return nil, err
	}
	cb := circuitbreaker.New(maxErrorPerSecond)

	return &Handler{
		httpClient: cli,
		retryCount: retryCount,
		breaker:    cb,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	body := ""
	code := 0
	for i := 0; i < h.retryCount; i += 1 {
		if !h.breaker.Allow() {
			body = "degradation"
			code = 206
			break
		}
		req, _ := http.NewRequest(http.MethodGet, "http://dependency:8080/do-work", nil)
		resp, err := h.httpClient.Do(req)
		if err != nil {
			body = "error"
			code = http.StatusInternalServerError
			h.breaker.ObserveErrors(1)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			body = "wrong status"
			code = http.StatusInternalServerError
			h.breaker.ObserveErrors(1)
			continue
		}

		body = "ok"
		code = 200
		break
	}

	w.WriteHeader(code)
	w.Write([]byte(body))
}
