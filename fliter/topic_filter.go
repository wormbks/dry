package filter

import (
	"strings"
	"sync"
)

// TopicFilter represents a filter for MQTT topics.
type TopicFilter struct {
	filterParts []string
	wildcards   map[int]bool
}

// NewTopicFilter creates a new TopicFilter instance based on the provided filter string.
func NewTopicFilter(filter string) *TopicFilter {
	parts := strings.Split(filter, "/")
	wildcards := make(map[int]bool)

	for i, part := range parts {
		if part == "+" || part == "#" {
			wildcards[i] = true
		}
	}

	return &TopicFilter{
		filterParts: parts,
		wildcards:   wildcards,
	}
}

// Matches checks if the topic matches the filter.
func (tf *TopicFilter) Matches(topic string) bool {
	parts := strings.Split(topic, "/")
	filterParts := tf.filterParts
	wildcards := tf.wildcards

	if len(parts) != len(filterParts) {
		return false
	}

	for i, filterPart := range filterParts {
		topicPart := parts[i]

		if !wildcards[i] && filterPart != topicPart {
			return false
		}
	}
	return true
}

// matchTopicFilters checks if the given topic matches any of the provided filters.
func MatchTopicFilters(filters []*TopicFilter, topic string) bool {
	for _, filter := range filters {
		if filter.Matches(topic) {
			return true
		}
	}
	return false
}

type TopicMatcher struct {
	filters []*TopicFilter
	cache   map[string]bool
	mutex   sync.Mutex
}

// NewTopicMatcher creates a new TopicMatcher instance based on the provided filters.
func NewTopicMatcher(filters []string) *TopicMatcher {
	tm := &TopicMatcher{
		filters: make([]*TopicFilter, 0),
		cache:   make(map[string]bool),
		mutex:   sync.Mutex{},
	}

	for _, filter := range filters {
		tm.AddFilter(filter)
	}

	return tm
}

// AddFilter adds a new filter to the TopicMatcher.
func (tm *TopicMatcher) AddFilter(filter string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.filters = append(tm.filters, NewTopicFilter(filter))
}

// RemoveFilter removes a filter from the TopicMatcher.
func (tm *TopicMatcher) RemoveFilter(filter string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	index := -1
	for i, tf := range tm.filters {
		if strings.Join(tf.filterParts, "/") == filter {
			index = i
			break
		}
	}

	if index != -1 {
		tm.filters = append(tm.filters[:index], tm.filters[index+1:]...)
	}
}

// Match checks if the given topic matches any of the filters in the TopicMatcher.
//
// It takes a string parameter named topic, which represents the topic to be matched.
// The function returns a boolean value indicating whether the topic matches any filter.
// If there are no filters in the TopicMatcher, the function returns true.
func (tm *TopicMatcher) Match(topic string) bool {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if len(tm.filters) == 0 { // empty filter
		return true
	}

	if matched, ok := tm.cache[topic]; ok {
		return matched
	}

	for _, filter := range tm.filters {
		if filter.Matches(topic) {
			tm.cache[topic] = true
			return true
		}
	}

	tm.cache[topic] = false
	return false
}
