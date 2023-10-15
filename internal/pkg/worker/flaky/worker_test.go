package flaky_test

import (
	"context"
	"sync"
	"testing"

	"load-testing/internal/pkg/worker/flaky"
)

func TestWorker_fast_path(t *testing.T) {
	w := flaky.New(1, 0, 0)
	err := w.Do(context.Background())
	if err != nil {
		t.Fatal()
	}
}

func TestWorker_flakiness(t *testing.T) {
	w := flaky.New(10, 100, 10)
	w.SetFlakiness(1)
	ctx := context.Background()
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := w.Do(ctx)
			if err == nil {
				t.Fatalf("worker should return error in case of flak")
			}
		}()
	}

	wg.Wait()
}

func TestWorker_concurrent_load(t *testing.T) {
	w := flaky.New(3, 100, 10)
	ctx := context.Background()
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := w.Do(ctx)
			if err != nil {
				t.Fatalf("worker should return nil in case of successful execution")
			}
		}()
	}

	wg.Wait()
}
