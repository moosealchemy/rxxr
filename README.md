# rxxr

Allows values to be published and subscribed to.
Subscribers will receive the most recently published value, if one exists.

Subscribed functions are run in their own goroutine.
Calling methods on the pipe are thread safe, but published values will not be locked.

### Constructor

```go
// Provide no options for default pipe.
pipe := rxxr.New[int]()
```

### Configure

```go
// To use the config:
pipe := rxxr.New[int](
    // Create a pipe already containing a value.
	// Will also imply SendOnSubscribe.
    UseInitialValue(42),
    // When a function is subscribed, call it with the current value if one exists.
	SendOnSubscribe,
	// Catch panics from subscribed functions.
    PanicHandler(func(r any) {
	    // Handle panic
    }),
)
```

### Subscribe

```go
pipe := rxxr.New[int]()

unsubscribe := pipe.Subscribe(func(i int) {
	// Do something
})

// Do something

unsubscribe()
```

### Publish

```go
pipe := rxxr.New[int]()
pipe.Subscribe(func (i int) {
	fmt.Printf("received: %d\n", i)
})
pipe.Publish(1)
pipe.Publish(2)
pipe.Publish(3)
pipe.Publish(5)
pipe.Publish(8)
pipe.Publish(13)

// received: 1
// received: 2
// received: 3
// received: 5
// received: 8
// received: 13
```

### Get Value

```go
pipe := rxxr.New[int]()
fmt.Println(pipe.Value())

pipe.Publish(42)
fmt.Println(pipe.Value())

// 0 false
// 42 true
```

### Close

```go
pipe := rxxr.New[int]()
pipe.Close()

// NOP
pipe.Publish(5)
```