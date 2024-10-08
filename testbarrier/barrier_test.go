package testbarrier_test

import (
	"artk.dev/testbarrier"
	"testing"
	"time"
)

func TestBarrier_succeeds_if_lifted_before_timeout_expires(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	go barrier.Lift()
	barrier.Wait(t, 100*365*24*time.Hour)
}

func TestBarrier_Wait_calls_FailNow_if_timeout_expires(t *testing.T) {
	t.Parallel()

	success := make(chan struct{})
	go func() {
		fakeT := &testingT{
			onHelper:  make(chan struct{}),
			onError:   make(chan struct{}),
			onFailNow: make(chan struct{}),
		}

		go func() {
			// Nothing can lift the barrier.
			// This guarantees that it will time out.
			barrier := testbarrier.New()
			barrier.Wait(fakeT, time.Nanosecond)
		}()

		<-fakeT.onHelper
		<-fakeT.onError
		<-fakeT.onFailNow
		success <- struct{}{}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	select {
	case <-success:
		// Hurrah!
	case <-ticker.C:
		t.Errorf("property was not satisfied within timeout")
	}
}

func TestBarrier_Wait_never_blocks_after_Lift(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	go barrier.Lift()

	for range 100 {
		barrier.Wait(t, 100*365*24*time.Hour)
	}
}
