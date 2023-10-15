package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	httpClient *http.Client
}

func FromEnv() (*Handler, error) {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT_MS"))
	if err != nil {
		return nil, err
	}

	cli := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	return &Handler{
		httpClient: cli,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(http.MethodGet, "http://dependency:8080/do-work", nil)
	resp, err := h.httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
		return
	}
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("wrong status"))
		return
	}

	w.Write([]byte("ok"))
}
