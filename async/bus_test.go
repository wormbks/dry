package async

import (
	"errors"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	assert.NotNil(t, bus, "Expected bus to be not nil")

	bus = NewMessageBus(0)

	assert.NotNil(t, bus, "Expected bus to be not nil")
}

func Test_Subscribe(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	err := bus.Subscribe("test", func() {})
	assert.NoError(t, err, "Expected no error when subscribing with a valid handler")

	err = bus.Subscribe("test", 2)
	assert.Error(t, err, "Expected an error when subscribing with an invalid handler")
}

func Test_Unsubscribe_InvalidHandler(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	err := bus.Subscribe("test", func(v bool) {})
	assert.NoError(t, err, "Expected no error when subscribing with a valid handler")

	err = bus.Unsubscribe("test", 2)
	assert.Error(t, err, "Expected an error when subscribing with an invalid handler")
}

func Test_Unsubscribe(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	handler01 := func() {}
	handler02 := func() {}

	err := bus.Subscribe("test-2subs", handler01)
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	err = bus.Subscribe("test-2subs", handler02)
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	err = bus.Unsubscribe("test-2subs", handler01)
	assert.NoError(t, err, "Expected no error when unsubscribing an existing handler")

	err = bus.Subscribe("test-1sub", handler01)
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	err = bus.Unsubscribe("test-1sub", handler01)
	assert.NoError(t, err, "Expected no error when unsubscribing an existing handler")

	err = bus.Unsubscribe("non-existed", func() {})
	assert.Error(t, err, "Expected an error when unsubscribing a non-existing topic")
}

func Test_Close(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	handler := func() {}

	err := bus.Subscribe("test", handler)
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	original, ok := bus.(*messageBus)
	assert.True(t, ok, "Could not cast message bus to its original type")

	assert.NotEmpty(t, original.handlers, "Expected handler to be subscribed to topic before closing")

	bus.Close("test")

	assert.Empty(t, original.handlers, "Expected handler to be unsubscribed from topic after closing")
}

func Test_Close_NoTopic(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	handler := func() {}

	err := bus.Subscribe("test", handler)
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	original, ok := bus.(*messageBus)
	assert.True(t, ok, "Could not cast message bus to its original type")

	assert.NotEmpty(t, original.handlers, "Expected handler to be subscribed to topic before closing")

	err = bus.Close("test-no-topic")
	assert.Error(t, err, "Expected an error when closing a non-existing topic")

	assert.NotEmpty(t, original.handlers, "Expected handler not to be unsubscribed from topic after closing")
}

func Test_Publish(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(2)

	first := false
	second := false

	err := bus.Subscribe("topic", func(v bool) {
		defer wg.Done()
		first = v
	})
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	err = bus.Subscribe("topic", func(v bool) {
		defer wg.Done()
		second = v
	})
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	bus.Publish("topic", true)

	wg.Wait()

	assert.True(t, first, "Expected first handler to be executed and set to true")
	assert.True(t, second, "Expected second handler to be executed and set to true")
}

func Test_Publish_No_Params(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	var wg sync.WaitGroup
	wg.Add(1)

	first := false

	err := bus.Subscribe("topic", func() {
		defer wg.Done()
		first = true
	})
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	err = bus.Publish("topic")
	assert.NoError(t, err, "Expected no error when publish")

	wg.Wait()

	assert.True(t, first, "Expected first handler to be executed and set to true")
}

func Test_Publish_NoHandler(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	err := bus.Subscribe("topic", func(v bool) {})
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	err = bus.Publish("topic-no-handler", true)
	assert.Error(t, err, "Expected an error when publishing a message without a handler")
}

func TestHandleError(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	err := bus.Subscribe("topic", func(out chan<- error) {
		out <- errors.New("throw error")
	})
	assert.NoError(t, err, "Expected no error when subscribing a valid handler")

	out := make(chan error)
	defer close(out)

	bus.Publish("topic", out)

	assert.Error(t, <-out, "Expected an error from the handler")
}
