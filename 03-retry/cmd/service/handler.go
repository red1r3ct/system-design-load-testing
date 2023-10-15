package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	httpClient *http.Client
	retryCount int
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

	return &Handler{
		httpClient: cli,
		retryCount: retryCount,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	body := ""
	code := 0
	for i := 0; i < h.retryCount; i += 1 {
		req, _ := http.NewRequest(http.MethodGet, "http://dependency:8080/do-work", nil)
		resp, err := h.httpClient.Do(req)
		if err != nil {
			body = "error"
			code = http.StatusInternalServerError
			continue
		}
		if resp.StatusCode != http.StatusOK {
			body = "wrong status"
			code = http.StatusInternalServerError
			continue
		}

		if i > 0 {
			fmt.Printf("saved in %d retries \n", i)
		}

		body = "ok"
		code = 200
		break
	}

	w.WriteHeader(code)
	w.Write([]byte(body))
}
