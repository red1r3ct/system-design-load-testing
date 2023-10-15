package flaky

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type FlakyWorker struct {
	// meanTime - mean time to do the work in ms
	meanTime int
	// stdTime - std of time to do the work in ms
	stdTime int
	// flakinessProba - number from 0 to 1 that represent probability of very long work
	flakinessProba float64

	pool chan struct{}
}

func New(poolSize int, mean int, std int) *FlakyWorker {
	pool := make(chan struct{}, poolSize)
	for i := 0; i < poolSize; i += 1 {
		pool <- struct{}{}
	}
	return &FlakyWorker{
		meanTime:       mean,
		stdTime:        std,
		flakinessProba: 0,
		pool:           pool,
	}
}

func (w *FlakyWorker) SetFlakiness(probability float64) {
	w.flakinessProba = probability
}

func (w *FlakyWorker) Do(ctx context.Context) error {
	if w.meanTime == 0 {
		// shortcut for fast response
		return nil
	}

	// read available resource or wait for it
	res := <-w.pool
	defer func() {
		// return resource the pool
		w.pool <- res
	}()
	// generate random time to finish the task
	timeToFinish := rand.NormFloat64()*float64(w.stdTime) + float64(w.meanTime)
	timeToFinish = math.Max(0, timeToFinish)
	workTime := time.Duration(timeToFinish) * time.Millisecond
	time.Sleep(workTime)

	if rand.Float64() < w.flakinessProba {
		// in case of flak, return err
		return fmt.Errorf("flaked")
	}

	return nil
}
