package async

import (
	"fmt"
	"sync"

	"github.com/go-faster/city"
)

type EventData struct {
	Data  any
	Topic string
}

type EventBus interface {
	Publish(topic string, data any) error
	Subscribe(topic string, ch EventChannel) uint64
	Unsubscribe(topic string, subscriptionID uint64)
}

// EventChannel is a channel which can accept a DataEvent
type EventChannel chan EventData

type subscriber struct {
	subscriptionID uint64
	ch             EventChannel
}

// eventBusImpl stores the information about subscribers interested for a particular topic
type eventBusImpl struct {
	subscribers map[string][]subscriber
	rm          sync.RWMutex
}

func NewEventBus() EventBus {
	return &eventBusImpl{
		subscribers: make(map[string][]subscriber),
	}
}

// Publish publishes the given data and topic to all subscribers.
// It locks read access to the subscribers map, defers unlocking,
// checks for subscribers for the topic, creates a dataEvent,
// ranges through the subscribers to send on their channels,
// and returns any error. If no subscribers are found, it returns
// ErrNoHandlerFound.
func (eb *eventBusImpl) Publish(topic string, data any) (err error) {
	eb.rm.RLock()
	defer eb.rm.RUnlock()
	if sbs, found := eb.subscribers[topic]; found {
		dataEvent := EventData{
			Data:  data,
			Topic: topic,
		}

		for _, sb := range sbs {
			select {
			case sb.ch <- dataEvent:
			default:
				// If the channel is full, drop the event.
				err = fmt.Errorf("event bus queue is full for topic %s", topic)
			}
		}
		return err
	}
	return ErrNoHandlerFound
}

// Subscribe registers a subscriber for a topic. It generates a unique
// subscription ID, adds the subscriber to the map of subscribers for
// that topic, and returns the subscription ID. It locks access to the
// subscribers map during this operation.
func (eb *eventBusImpl) Subscribe(topic string, ch EventChannel) uint64 {
	eb.rm.Lock()
	defer eb.rm.Unlock()
	// Generate a unique subscription ID
	subscriptionID := generateUInt64ID(topic, len(eb.subscribers[topic])+1)
	s := subscriber{subscriptionID, ch}

	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, s)
	} else {
		sbs := make([]subscriber, 0, 5)
		sbs = append(sbs, s)
		eb.subscribers[topic] = sbs
	}

	return subscriptionID
}

// Unsubscribe removes the subscriber with the given subscription ID
// from the subscribers list for the given topic. It locks access to
// the subscribers map during the operation.
func (eb *eventBusImpl) Unsubscribe(topic string, subscriptionID uint64) {
	eb.rm.Lock()
	defer eb.rm.Unlock()
	if sbs, found := eb.subscribers[topic]; found {
		for i, sb := range sbs {
			if sb.subscriptionID == subscriptionID {
				eb.subscribers[topic] = append(sbs[:i], sbs[i+1:]...)
				break
			}
		}
	}
}

// generateUInt64ID generates a unique 64-bit unsigned integer ID
// by combining the given string and integer. It hashes the
// concatenated byte slice using the CityHash64 hashing function.
func generateUInt64ID(str string, num int) uint64 {
	// Combine string and int64 into a single byte slice
	data := []byte(str)
	data = append(data, []byte(fmt.Sprintf("%d", num))...)

	// Hash the data using CityHash64
	id := city.Hash64(data)

	return id
}
