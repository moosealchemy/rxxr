package rxxr

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("new with nil config", func(t *testing.T) {
		p := New[any](nil)
		if p == nil {
			t.Fail()
		}
	})

	t.Run("new with config", func(t *testing.T) {
		p := New[any]()
		if p == nil {
			t.Fail()
		}
	})

	t.Run("new with initial value", func(t *testing.T) {
		p := New[int](UseInitialValue(5))
		v, ok := p.Value()
		if p == nil || v != 5 || !ok {
			t.Fail()
		}
	})

	t.Run("new with send on subscribe", func(t *testing.T) {
		p := New[int](SendOnSubscribe[int])
		if p == nil {
			t.Fail()
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		p.Publish(5)
		p.Subscribe(func(i int) {
			defer cancel()
			if i != 5 {
				t.Fail()
			}
		})
		watchWithTimeout(t, ctx.Done())
	})

	t.Run("new with panic handler", func(t *testing.T) {
		var (
			tv          = 5
			ctx, cancel = context.WithCancel(context.Background())
			p           = New[int](
				UseInitialValue(tv),
				PanicHandler[int](func(r any) {
					defer cancel()
					if i, ok := r.(int); !ok || i != tv {
						t.Fail()
					}
				}))
		)
		defer cancel()

		if p == nil {
			t.Fail()
			return
		}
		p.Subscribe(func(i int) { panic(i) })
		watchWithTimeout(t, ctx.Done())
	})
}
