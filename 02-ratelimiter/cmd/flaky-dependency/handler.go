package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"load-testing/internal/pkg/worker/flaky"
)

type Handler struct {
	w *flaky.FlakyWorker
}

func FromEnv() (*Handler, error) {
	poolSize, err := strconv.Atoi(os.Getenv("POOL_SIZE"))
	if err != nil {
		return nil, err
	}
	mean, err := strconv.Atoi(os.Getenv("MEAN_WORK_MS"))
	if err != nil {
		return nil, err
	}
	std, err := strconv.Atoi(os.Getenv("STD_WORK_MS"))
	if err != nil {
		return nil, err
	}
	flakiness, err := strconv.ParseFloat(os.Getenv("DEFAULT_FLAKINESS"), 64)
	if err != nil {
		return nil, err
	}

	w := flaky.New(poolSize, mean, std)
	w.SetFlakiness(flakiness)

	return &Handler{
		w: w,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	err := h.w.Do(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	} else {
		w.Write([]byte("ok"))
	}
}

type SetFlakinessRequest struct {
	Flakiness float64 `json:"flakiness"`
}

func (h *Handler) SetFlakiness(w http.ResponseWriter, r *http.Request) {
	req := &SetFlakinessRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil || req.Flakiness < 0 || req.Flakiness > 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error"))
		return
	}

	h.w.SetFlakiness(req.Flakiness)
	w.Write([]byte("ok"))
}
