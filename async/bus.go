package async

// It is based on https://github.com/vardius/message-bus
// Copyright (c) 2017-present Rafał Lorenz

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

const DefHandlerQueueSize = 64

var (
	ErrNoHandlerFound = errors.New("no bus handler found")
	ErrTopicNotFound  = errors.New("bus topic not found")
	ErrQueueFull      = errors.New("bus queue is full")
)

// MessageBus implements publish/subscribe messaging paradigm
type MessageBus interface {
	// Publish publishes arguments to the given topic subscribers
	// Publish block only when the buffer of one of the subscribers is full.
	Publish(topic string, args ...interface{}) error
	// Close unsubscribe all handlers from given topic
	Close(topic string) error
	// Subscribe subscribes to the given topic
	Subscribe(topic string, fn interface{}) error
	// Unsubscribe unsubscribe handler from the given topic
	Unsubscribe(topic string, fn interface{}) error
}

type handlersMap map[string][]*msgHandler

type msgHandler struct {
	callback reflect.Value
	queue    chan []reflect.Value
}

type messageBus struct {
	handlerQueueSize int
	mtx              sync.RWMutex
	handlers         handlersMap
}

// Publish publishes a message to the given topic in the message bus.
//
// It takes a topic string and a variable number of arguments as its parameters.
// The function returns an error.
func (b *messageBus) Publish(topic string, args ...interface{}) (err error) {
	rArgs := buildHandlerArgs(args)

	b.mtx.RLock()
	defer b.mtx.RUnlock()

	if hs, ok := b.handlers[topic]; ok {
		for _, h := range hs {
			h.queue <- rArgs
		}
	} else {
		err = ErrNoHandlerFound
	}
	return err
}

// Subscribe subscribes to a topic and registers a callback function to be executed when a message is received.
//
// Parameters:
// - topic: the topic to subscribe to (string).
// - fn: the callback function to be executed (interface{}).
//
// Returns:
// - error: if there is an error validating the callback function.
func (b *messageBus) Subscribe(topic string, fn interface{}) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	h := &msgHandler{
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value, b.handlerQueueSize),
	}

	go func() {
		for args := range h.queue {
			h.callback.Call(args)
		}
	}()

	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.handlers[topic] = append(b.handlers[topic], h)

	return nil
}

// Unsubscribe unsubscribes a handler function from a specific topic in the message bus.
//
// It takes in the topic string and the handler function fn as parameters.
// The handler function fn must be a valid handler, otherwise an error is returned.
// It returns an error if the topic is not found in the message bus.
func (b *messageBus) Unsubscribe(topic string, fn interface{}) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	rv := reflect.ValueOf(fn)

	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.handlers[topic]; ok {
		for i, h := range b.handlers[topic] {
			if h.callback == rv {
				close(h.queue)

				if len(b.handlers[topic]) == 1 {
					delete(b.handlers, topic)
				} else {
					b.handlers[topic] = append(b.handlers[topic][:i], b.handlers[topic][i+1:]...)
				}
			}
		}

		return nil
	}

	return ErrTopicNotFound
}

// Close closes the message bus for a given topic.
//
// It takes a string parameter `topic` which represents the topic to be closed.
// The function returns an error indicating if the topic was not found.
func (b *messageBus) Close(topic string) (err error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.handlers[topic]; ok {
		for _, h := range b.handlers[topic] {
			close(h.queue)
		}

		delete(b.handlers, topic)
	} else {
		err = ErrTopicNotFound
	}

	return err
}

// isValidHandler checks if the given function is a valid handler.
//
// fn: the function to be checked.
// err: an error indicating whether the function is valid or not.
// Returns: nil if the function is valid, otherwise an error indicating the reason.
func isValidHandler(fn interface{}) (err error) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		err = fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	return err
}

// buildHandlerArgs creates an array of reflect.Value objects from an array of interface{} objects.
//
// args: The array of interface{} objects to be converted.
// []reflect.Value: The array of reflect.Value objects created from the input array.
func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}

// NewMessageBus creates new MessageBus
// handlerQueueSize sets buffered channel length per subscriber
func NewMessageBus(handlerQueueSize int) MessageBus {
	if handlerQueueSize < 1 {
		handlerQueueSize = DefHandlerQueueSize
	}

	return &messageBus{
		handlerQueueSize: handlerQueueSize,
		handlers:         make(handlersMap),
	}
}
