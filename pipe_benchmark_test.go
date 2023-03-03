package rxxr

import (
	"math"
	"sync"
	"testing"
)

func BenchmarkPipe_Publish(b *testing.B) {
	p := New[int](PanicHandler[int](func(r any) {
		b.Fail()
		b.Log("benchmark panicked")
	}))

	subscribers := math.MaxInt16
	wg := new(sync.WaitGroup)
	for i := 0; i < subscribers; i++ {
		p.Subscribe(func(int) { wg.Done() })
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(subscribers)
		p.Publish(i)
		wg.Wait()
		wg = new(sync.WaitGroup)
	}
}
func BenchmarkDone(b *testing.B) {
	subscribers := math.MaxInt16
	wg := new(sync.WaitGroup)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(subscribers)
		for i := 0; i < subscribers; i++ {
			wg.Done()
		}
		wg.Wait()
		wg = new(sync.WaitGroup)
	}
}
