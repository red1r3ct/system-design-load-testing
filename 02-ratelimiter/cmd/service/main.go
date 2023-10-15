package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/time/rate"
)

func ratelimit(limitRPS int, next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(limitRPS), 3)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if limiter.Allow() == false {
            http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

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
	http.Handle("/call", ratelimit(limitRPS, http.HandlerFunc(handler.Handle)))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		os.Exit(3)
	}
}
