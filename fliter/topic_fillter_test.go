package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicFilter_Matches(t *testing.T) {
	tf := NewTopicFilter("sensors/+/temperature")

	// Matching topics
	assert.True(t, tf.Matches("sensors/1/temperature"))
	assert.True(t, tf.Matches("sensors/2/temperature"))
	assert.True(t, tf.Matches("sensors/3/temperature"))

	// Non-matching topics
	assert.False(t, tf.Matches("sensors/1/humidity"))
	assert.False(t, tf.Matches("sensors/2/humidity"))
	assert.False(t, tf.Matches("devices/1/temperature"))
}

func TestMatchTopicFilters(t *testing.T) {
	filters := []*TopicFilter{
		NewTopicFilter("sensors/+/temperature"),
		NewTopicFilter("sensors/2/#"),
		NewTopicFilter("devices/+/status"),
	}

	// Matching topics
	assert.True(t, MatchTopicFilters(filters, "sensors/1/temperature"))
	assert.True(t, MatchTopicFilters(filters, "sensors/2/temperature"))
	assert.True(t, MatchTopicFilters(filters, "devices/office/status"))

	// Non-matching topics
	assert.False(t, MatchTopicFilters(filters, "sensors/1/humidity"))
	assert.False(t, MatchTopicFilters(filters, "devices/kitchen/temperature"))
}
