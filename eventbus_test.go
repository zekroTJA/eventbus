package eventbus

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPubSubAsync(t *testing.T) {
	const ceiling = 1000

	b := New[uint32]()

	var sumA, sumB uint32
	var wgPub, wgSub sync.WaitGroup
	wgPub.Add(2)
	wgSub.Add(2)

	cA, unsubA := b.Subscribe()
	cB, unsubB := b.Subscribe()

	go func() {
		for v := range cA {
			atomic.AddUint32(&sumA, v)
		}
		wgSub.Done()
	}()

	go func() {
		for v := range cB {
			atomic.AddUint32(&sumB, v)
		}
		wgSub.Done()
	}()

	go func() {
		for i := 2; i <= ceiling; i += 2 {
			b.Publish(uint32(i))
		}
		wgPub.Done()
	}()

	go func() {
		for i := 1; i <= ceiling; i += 2 {
			b.Publish(uint32(i))
		}
		wgPub.Done()
	}()

	wgPub.Wait()

	unsubA()
	unsubB()

	wgSub.Wait()

	expectSumA := sum(ceiling)
	if sumA != expectSumA {
		t.Errorf("subA was %d != expected %d", sumA, expectSumA)
	}

	expectSumB := expectSumA
	if sumB != expectSumB {
		t.Errorf("subA was %d != expected %d", sumB, expectSumA)
	}
}

func TestPubSubFuncAsync(t *testing.T) {
	const ceiling = 1000

	b := New[uint32]()

	var sumA, sumB uint32
	var wgPub sync.WaitGroup
	wgPub.Add(2)

	b.SubscribeFunc(func(v uint32) {
		atomic.AddUint32(&sumA, v)
	})
	b.SubscribeFunc(func(v uint32) {
		atomic.AddUint32(&sumB, v)
	})

	go func() {
		for i := 2; i <= ceiling; i += 2 {
			b.Publish(uint32(i))
		}
		wgPub.Done()
	}()

	go func() {
		for i := 1; i <= ceiling; i += 2 {
			b.Publish(uint32(i))
		}
		wgPub.Done()
	}()

	wgPub.Wait()

	// Because we don't track subscriptions in channels,
	// we need to wait here a fair bit of time until all
	// subscription go routines have finished.
	time.Sleep(100 * time.Millisecond)

	expectSumA := sum(ceiling)
	if sumA != expectSumA {
		t.Errorf("subA was %d != expected %d", sumA, expectSumA)
	}

	expectSumB := expectSumA
	if sumB != expectSumB {
		t.Errorf("subA was %d != expected %d", sumB, expectSumA)
	}
}

func TestSubscribeFuncParams(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("passing nil did not panic")
		}
	}()

	b := New[any]()
	b.SubscribeFunc(nil)
}

// --- helpers ---

func sum(ceiling uint32) uint32 {
	var r, i uint32 = 0, 1
	for ; i <= ceiling; i++ {
		r += i
	}
	return r
}
