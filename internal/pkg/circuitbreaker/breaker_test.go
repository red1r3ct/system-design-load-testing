package circuitbreaker_test

import (
	"testing"
	"time"

	"load-testing/internal/pkg/circuitbreaker"
)

func TestCircuitBreaker_closed_on_high_rate(t *testing.T) {
	cb := circuitbreaker.New(10)

	if !cb.Allow() {
		t.Fatalf("circuit breaker should open by default")
	}

	cb.ObserveErrors(20)
	if cb.Allow() {
		t.Fatalf("circuit breaker closed under high error rate")
	}

	time.Sleep(3 * time.Second)
	if !cb.Allow() {
		t.Fatalf("circuit breaker should be released after enough time")
	}
}
