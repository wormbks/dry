package iter

// LoopIteratorr represents an endless/loop iterator for a generic array.
type LoopIteratorr[T any] struct {
	data  []T
	index int
}

// NewIterator creates a new iterator for the given array.
func NewLoopIterator[T any](data []T) *LoopIteratorr[T] {
	return &LoopIteratorr[T]{data: data, index: 0}
}

// Next returns the next value from the iterator.
func (iter *LoopIteratorr[T]) Next() T {
	val := iter.data[iter.index]
	iter.index = (iter.index + 1) % len(iter.data) // Move to the next index, loop back if necessary
	return val
}

func (iter *LoopIteratorr[T]) Reset() {
	iter.index = 0
}

// IntRangeIterator is a type for iterating over a range of integers.
type IntRangeIterator struct {
	Start, End, Current int
}

// NewIntRangeIterator creates a new IntRangeIterator instance.
func NewIntRangeIterator(start, end int) *IntRangeIterator {
	return &IntRangeIterator{Start: start, End: end, Current: start - 1}
}

// Next returns the next integer in the range and a boolean indicating if there are more values.
func (it *IntRangeIterator) Next() (int, bool) {
	it.Current++
	return it.Current, it.Current <= it.End
}

// Reset resets the IntRangeIterator to its initial state.
func (it *IntRangeIterator) Reset() {
	it.Current = it.Start - 1
}

// LoopingIntRangeIterator is a type for iterating over a range of integers in a loop.
type LoopingIntRangeIterator struct {
	Start, End, Current int
}

// NewLoopingIntRangeIterator creates a new LoopingIntRangeIterator instance.
func NewLoopingIntRangeIterator(start, end int) *LoopingIntRangeIterator {
	return &LoopingIntRangeIterator{Start: start, End: end, Current: start - 1}
}

// Next returns the next integer in the range and loops back to the start if it reaches the end.
func (it *LoopingIntRangeIterator) Next() int {
	it.Current += 1
	if it.Current > it.End {
		it.Current = it.Start
	}
	return it.Current
}

// Reset resets the LoopingIntRangeIterator.
func (it *LoopingIntRangeIterator) Reset() {
	it.Current = it.Start - 1
}
