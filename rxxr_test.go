package rxxr

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"testing"
)

func watchWithTimeout(t *testing.T, done <-chan struct{}) {
	t.Helper()
	ctx := context.Background()
	if deadline, ok := t.Deadline(); ok {
		var cancel func()
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}

	select {
	case <-done:
	case <-ctx.Done():
		t.Fail()
	}
}

func Test_value_Publish(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var (
			p           = New[int]()
			i           = rand.Int()
			ctx, cancel = context.WithCancel(context.Background())
		)

		p.Subscribe(func(v int) {
			defer cancel()
			if v != i {
				t.Logf("expected %d, got %d", i, v)
				t.Fail()
			}
		})
		p.Publish(i)
		watchWithTimeout(t, ctx.Done())
	})

	t.Run("publish many subscribers", func(t *testing.T) {
		p := New[int](nil)
		wg := new(sync.WaitGroup)
		count := math.MaxInt16
		ctx, cancel := context.WithCancel(context.Background())

		for i := 0; i <= count; i++ {
			wg.Add(1)
			p.Subscribe(func(v int) {
				defer wg.Done()
			})
		}

		go func() {
			t.Helper()
			defer cancel()
			wg.Wait()
		}()

		p.Publish(0)
		watchWithTimeout(t, ctx.Done())
	})
}

func Test_value_Subscribe(t *testing.T) {
	t.Run("basic subscription", func(t *testing.T) {
		p := New[any](nil)
		s := p.Subscribe(func(any) {})
		if s == nil {
			t.Fail()
		}
		s()
	})
}

func Test_value_Value(t *testing.T) {
	t.Run("no value", func(t *testing.T) {
		p := New[any](nil)
		if _, ok := p.Value(); ok {
			t.Fail()
		}
	})

	t.Run("has value", func(t *testing.T) {
		i := rand.Int()
		p := New[int](UseInitialValue(i))
		if v, ok := p.Value(); !ok {
			t.Fail()
		} else if v != i {
			t.Fail()
		}
	})
}
