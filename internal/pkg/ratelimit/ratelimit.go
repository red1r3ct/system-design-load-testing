package ratelimit

import (
	"net/http"

	"golang.org/x/time/rate"
)

func Middleware(limitRPS int, next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(limitRPS), 3)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if limiter.Allow() == false {
            http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}