package str

import (
	"strings"
	"sync"
)

const EmptyString = "__empty__"

// StringReplacer is a generic type to replace specified substrings
// in strings and cache the modified strings.
type StringReplacer struct {
	cache               sync.Map          // Cache for storing modified strings, using string keys
	replacePairs        map[string]string // ReplacePairs is a map of strings to replace with their corresponding replacement strings.
	reverseReplacePairs map[string]string // ReverseReplacePairs is a reverse mapping of replacement strings to their original values.
}

// NewStringReplacer creates a new instance of StringReplacer with specified replaceable and replacement strings.
func NewStringReplacer(replacePairs map[string]string) *StringReplacer {
	// Create a reverse mapping for bidirectional transformation
	reverseReplacePairs := make(map[string]string)
	for k, v := range replacePairs {
		reverseReplacePairs[v] = k
	}

	return &StringReplacer{
		replacePairs:        replacePairs,
		reverseReplacePairs: reverseReplacePairs,
	}
}

// Replace replaces the specified substrings in the input string with their corresponding values.
func (sr *StringReplacer) Replace(input string) string {
	// Create a key for the cache
	cacheKey := input

	// Check if the modified string is already in the cache
	if cachedValue, ok := sr.getFromCache(cacheKey); ok {
		return cachedValue
	}

	// Replace specified substrings in the input string
	modifiedString := input
	for replace, replacement := range sr.replacePairs {
		modifiedString = strings.Replace(modifiedString, replace, replacement, -1)
	}

	// Cache the modified string
	sr.addToCache(cacheKey, modifiedString)

	return modifiedString
}

// ReverseReplace replaces the specified substrings in the input string with their corresponding original values.
func (sr *StringReplacer) ReverseReplace(input string) string {
	// Create a key for the cache
	cacheKey := input

	// Check if the modified string is already in the cache
	if cachedValue, ok := sr.getFromCache(cacheKey); ok {
		return cachedValue
	}

	// Replace specified substrings in the input string with their original values
	modifiedString := input
	for replacement, original := range sr.reverseReplacePairs {
		modifiedString = strings.Replace(modifiedString, replacement, original, -1)
	}

	// Cache the modified string
	sr.addToCache(cacheKey, modifiedString)

	return modifiedString
}

// addToCache adds the modified string to the cache.
func (sr *StringReplacer) addToCache(key, value string) {
	sr.cache.Store(key, value)
}

// getFromCache retrieves the modified string from the cache.
func (sr *StringReplacer) getFromCache(key string) (string, bool) {
	if cachedValue, ok := sr.cache.Load(key); ok {
		return cachedValue.(string), true
	}
	return "", false
}
