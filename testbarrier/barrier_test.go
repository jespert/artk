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
			barrier := testbarrier.New()
			barrier.Wait(fakeT, time.Nanosecond)
		}()

		<-fakeT.onHelper
		<-fakeT.onError
		<-fakeT.onFailNow
		success <- struct{}{}
	}()

	select {
	case <-success:
		// Hurrah!
	case <-time.NewTicker(5 * time.Second).C:
		t.Errorf("property was not satisfied within timeout")
	}
}
