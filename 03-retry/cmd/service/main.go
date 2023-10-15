package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"load-testing/internal/pkg/ratelimit"
)

func main() {
	limitRPS, err := strconv.Atoi(os.Getenv("RATELIMIT_RPS"))
	if err != nil {
		fmt.Println("Rate limit is not specified")
		os.Exit(4)
	}

	handler, err := FromEnv()
	if err != nil {
		fmt.Println("failed to load handler from env")
		os.Exit(2)
	}
	http.Handle("/call", ratelimit.Middleware(limitRPS, http.HandlerFunc(handler.Handle)))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		os.Exit(3)
	}
}
