package rxxr

type Option[T any] func(*Pipe[T])

// UseInitialValue Initializes a pipe with the provided value.
func UseInitialValue[T any](t T) Option[T] {
	return func(p *Pipe[T]) {
		SendOnSubscribe(p)
		p.value = &t
	}
}

// SendOnSubscribe sets if the current value should be sent to new subscribers,
// if a value is not sent the subscriber will receive no values.
func SendOnSubscribe[T any](p *Pipe[T]) { p.sendOnSubscribe = true }

// PanicHandler sets the function to receive panics from running subscribed funcs
func PanicHandler[T any](fn func(any)) Option[T] { return func(p *Pipe[T]) { p.panicHandler = fn } }
