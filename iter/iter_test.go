package iter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	// Test NewIterator
	data := []int{1, 2, 3, 4, 5}
	iterator := NewLoopIterator(data)

	// Test Next method for each value in the array
	for i := 0; i < len(data); i++ {
		expectedValue := data[i]
		value := iterator.Next()
		assert.Equal(t, expectedValue, value, "Unexpected value from Iterator")
	}

	// Test Next method after reaching the end of the array
	value := iterator.Next()
	assert.Equal(t, data[0], value, "Unexpected value from Iterator after reaching the end")
}

func TestIteratorAfterReset(t *testing.T) {
	// Test NewIterator
	data := []int{1, 2, 3, 4, 5}
	iterator := NewLoopIterator(data)

	// Test Next method for each value in the array
	for i := 0; i < len(data); i++ {
		expectedValue := data[i]
		value := iterator.Next()
		assert.Equal(t, expectedValue, value, "Unexpected value from Iterator")
	}

	// Reset the iterator
	iterator.Reset()

	// Test Next method after resetting
	for i := 0; i < len(data); i++ {
		expectedValue := data[i]
		value := iterator.Next()
		assert.Equal(t, expectedValue, value, "Unexpected value from Iterator after reset")
	}
}

func TestIntRangeIterator(t *testing.T) {
	iterator := NewIntRangeIterator(1, 5)

	// Test Next method for each value in the range
	for i := 1; i <= 5; i++ {
		value, hasNext := iterator.Next()
		assert.True(t, hasNext, "Expected hasNext to be true")
		assert.Equal(t, i, value, "Unexpected value from IntRangeIterator")
	}

	// Test Next method after reaching the end of the range
	value, hasNext := iterator.Next()
	assert.False(t, hasNext, "Expected hasNext to be false")
	assert.Equal(t, 6, value, "Unexpected value from IntRangeIterator after reaching the end")
}

func TestLoopingIntRangeIterator(t *testing.T) {
	iterator := NewLoopingIntRangeIterator(1, 5)

	// Test Next method for each value in the range
	val := []int{0, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	for i := 1; i <= 10; i++ {
		value := iterator.Next()
		expectedValue := val[i]
		assert.Equal(t, expectedValue, value, "Unexpected value from LoopingIntRangeIterator")
	}
}

func TestLoopingIntRangeIteratorAfterReset(t *testing.T) {
	iterator := NewLoopingIntRangeIterator(1, 5)
	// Test Next method for each value in the range
	val := []int{0, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5}

	// Test Next method for each value in the range
	for i := 1; i <= 5; i++ {
		expectedValue := val[i]
		value := iterator.Next()
		assert.Equal(t, expectedValue, value, "Unexpected value from LoopingIntRangeIterator")
	}

	// Reset the iterator
	iterator.Reset()

	// Test Next method after resetting
	for i := 1; i <= 10; i++ {
		expectedValue := val[i]
		value := iterator.Next()
		assert.Equal(t, expectedValue, value, "Unexpected value from LoopingIntRangeIterator after reset")
	}
}
