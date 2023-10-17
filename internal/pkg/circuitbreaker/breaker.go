package circuitbreaker

import (
	"time"

	"golang.org/x/time/rate"
)

type CircuitBreaker struct {
	maxErrorPerSecond float64
	limiter           *rate.Limiter
}

func New(maxErrorPerSecond float64) *CircuitBreaker {
	return &CircuitBreaker{
		maxErrorPerSecond: maxErrorPerSecond,
		limiter:           rate.NewLimiter(rate.Limit(maxErrorPerSecond), 1),
	}
}

func (cb *CircuitBreaker) ObserveErrors(n int) {
	cb.limiter.ReserveN(time.Now(), n)
}

func (cb *CircuitBreaker) Allow() bool {
	return cb.limiter.Tokens() > 0
}
