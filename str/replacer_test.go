package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringReplacerReplace(t *testing.T) {
	// Create a StringReplacer instance with replaceable and replacement strings
	replacePairs := map[string]string{
		"IP":       "newIPAddress",
		"target":   "newTarget",
		"original": "replacement",
	}
	replacer := NewStringReplacer(replacePairs)

	// Test case 1
	input := "QLM/IP/target/original"
	expectedResult := "QLM/newIPAddress/newTarget/replacement"
	result := replacer.Replace(input)
	assert.Equal(t, expectedResult, result, "Unexpected result")

	// Check if the cached value matches the expected value
	cachedValue, ok := replacer.getFromCache(input)
	assert.True(t, ok, "Unexpected cache status")
	assert.Equal(t, expectedResult, cachedValue, "Unexpected cached value")

	// Test case 2
	input = "QLM/IP/target/original"
	expectedResult = "QLM/newIPAddress/newTarget/replacement"
	result = replacer.Replace(input)
	assert.Equal(t, expectedResult, result, "Unexpected result")

	// Check if the cached value matches the expected value
	cachedValue, ok = replacer.getFromCache(input)
	assert.True(t, ok, "Unexpected cache status")
	assert.Equal(t, "QLM/newIPAddress/newTarget/replacement", cachedValue, "Unexpected cached value")

	// Test case 3
	// Check if the cached value matches the expected value
	input = "Different/Input"
	expectedResult = "Different/Input"
	cachedValue, ok = replacer.getFromCache(input)
	assert.False(t, ok, "Unexpected cache status")
	assert.Equal(t, "", cachedValue, "Unexpected cached value")

	result = replacer.Replace(input)
	assert.Equal(t, expectedResult, result, "Unexpected result")
}

func TestStringReplacerReverseReplace(t *testing.T) {
	// Create a StringReplacer instance with replaceable and replacement strings
	replacePairs := map[string]string{
		"IP":       "newIPAddress",
		"target":   "newTarget",
		"original": "replacement",
	}
	replacer := NewStringReplacer(replacePairs)

	// Test case 1
	input := "QLM/newIPAddress/newTarget/replacement"
	expectedResult := "QLM/IP/target/original"

	// Check if the cached value matches the expected value
	cachedValue, ok := replacer.getFromCache(input)
	assert.False(t, ok, "Unexpected cache status")
	assert.Equal(t, "", cachedValue, "Unexpected cached value")

	result := replacer.ReverseReplace(input)
	assert.Equal(t, expectedResult, result, "Unexpected result")

	// Test case 2
	// Check if the cached value matches the expected value
	cachedValue, ok = replacer.getFromCache(input)
	assert.True(t, ok, "Unexpected cache status")
	assert.Equal(t, "QLM/IP/target/original", cachedValue, "Unexpected cached value")

	// Test case 3
	input = "Different/Input"
	expectedResult = "Different/Input"
	result = replacer.ReverseReplace(input)
	assert.Equal(t, expectedResult, result, "Unexpected result")
}
