package async

import (
	"testing"
)

func TestEventBus_Publish(t *testing.T) {
	eb := NewEventBus()
	ch := make(EventChannel, 10)
	topic := "testTopic"
	expectedData := "someData"

	// Subscribe to the topic
	subscriptionID := eb.Subscribe(topic, ch)

	// Publish an event with the data
	err := eb.Publish(topic, expectedData)
	if err != nil {
		t.Errorf("Publish failed: %v", err)
	}

	// Receive the event from the channel
	event := <-ch

	// Assert that the received data matches the expected data
	if event.Data != expectedData {
		t.Errorf("Received data does not match expected data: %v != %v", event.Data, expectedData)
	}

	// Unsubscribe from the topic
	eb.Unsubscribe(topic, subscriptionID)
}

func TestEventBus_Subscribe(t *testing.T) {
	eb := NewEventBus()
	topic := "testTopic"
	ch := make(EventChannel, 10)

	// Subscribe to the topic
	subscriptionID1 := eb.Subscribe(topic, ch)
	subscriptionID2 := eb.Subscribe(topic, ch) // Subscribe again with the same channel

	// Assert that different subscription IDs are generated
	if subscriptionID1 == subscriptionID2 {
		t.Errorf("Expected different subscription IDs, got: %d", subscriptionID1)
	}

	// Unsubscribe from the topic
	eb.Unsubscribe(topic, subscriptionID1)
	eb.Unsubscribe(topic, subscriptionID2)
}

func TestEventBus_Unsubscribe(t *testing.T) {
	eb := NewEventBus()
	topic := "testTopic"
	ch := make(EventChannel, 10)

	// Subscribe to the topic
	subscriptionID := eb.Subscribe(topic, ch)

	// Unsubscribe from the topic using the subscription ID
	eb.Unsubscribe(topic, subscriptionID)

	// Publish an event and try to receive it
	eb.Publish(topic, "someData")
	select {
	case <-ch:
		t.Error("Received event on unsubscribed channel")
	default:
		// Expected behavior
	}
}

func TestGenerateUInt64ID(t *testing.T) {
	str := "testString"
	num := 123

	// Generate the ID
	id1 := generateUInt64ID(str, num)
	id2 := generateUInt64ID(str, num) // Generate again with the same input

	// Assert that the IDs are the same for the same input
	if id1 != id2 {
		t.Errorf("Generated IDs are not the same for the same input: %d != %d", id1, id2)
	}
}
