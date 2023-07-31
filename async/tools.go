package async

import (
	"errors"
	"io"
	"sync"
	"time"
)

// DelayedWriter is a custom writer that delays the write operation.
type DelayedWriter struct {
	writer io.Writer
	delay  time.Duration
}

// NewDelayedWriter creates a new DelayedWriter with the provided writer and delay.
//
// writer: the io.Writer to write to.
// delay: the time.Duration to delay the write operation.
// Returns a pointer to the created DelayedWriter.
func NewDelayedWriter(writer io.Writer, delay time.Duration) *DelayedWriter {
	return &DelayedWriter{
		writer: writer,
		delay:  delay,
	}
}

// Write method introduces the delay using time.Sleep and then delegates
// the actual writing to the embedded io.Writer.
func (dw *DelayedWriter) Write(data []byte) (int, error) {
	time.Sleep(dw.delay) // Introduce the delay before writing
	return dw.writer.Write(data)
}

// Close is a method of the DelayedWriter struct that closes the writer.
//
// It does not take any parameters.
// It returns an error.
func (dw *DelayedWriter) Close() error {
	return nil
}

// ErrorAfterWriter is a custom writer that returns an error after
// a specified number of write operations.
type ErrorAfterWriter struct {
	writer         io.Writer
	maxWriteCount  int
	currentWrites  int
	errAfterWrites int
	closed         bool
	mutex          sync.Mutex
}

// NewErrorAfterWriter creates a new ErrorAfterWriter.
//
// It takes in an io.Writer, maxWriteCount, and errAfterWrites as parameters.
// Returns a pointer to an ErrorAfterWriter.
func NewErrorAfterWriter(writer io.Writer, maxWriteCount, errAfterWrites int) *ErrorAfterWriter {
	return &ErrorAfterWriter{
		writer:         writer,
		maxWriteCount:  maxWriteCount,
		currentWrites:  0,
		errAfterWrites: errAfterWrites,
		closed:         false,
	}
}

// Write method keeps track of the number of write operations (currentWrites)
// and returns an error after reaching the specified limit (errAfterWrites).
func (ew *ErrorAfterWriter) Write(data []byte) (int, error) {
	ew.mutex.Lock()
	defer ew.mutex.Unlock()

	if ew.closed {
		return 0, errors.New("writer closed")
	}

	if ew.currentWrites >= ew.errAfterWrites {
		return 0, errors.New("error after write operations")
	}

	ew.currentWrites++
	return ew.writer.Write(data)
}

// Close closes the ErrorAfterWriter. It sets the 'closed' field to true.
// It returns an error if any.
func (ew *ErrorAfterWriter) Close() error {
	ew.mutex.Lock()
	defer ew.mutex.Unlock()

	ew.closed = true
	return nil
}
