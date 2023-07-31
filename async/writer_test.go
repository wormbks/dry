package async

// async_test.go

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// BufferWriteCloser is a wrapper around bytes.Buffer that implements the WriteCloser interface.
type BufferWriteCloser struct {
	buf *bytes.Buffer
}

// NewBufferWriteCloser creates a new BufferWriteCloser.
func NewBufferWriteCloser() *BufferWriteCloser {
	return &BufferWriteCloser{
		buf: &bytes.Buffer{},
	}
}

// Write appends data to the buffer.
func (bwc *BufferWriteCloser) Write(p []byte) (n int, err error) {
	return bwc.buf.Write(p)
}

// Close does nothing since bytes.Buffer doesn't require explicit closing.
func (bwc *BufferWriteCloser) Close() error {
	return nil
}

// String returns the string representation of the BufferWriteCloser.
func (bwc *BufferWriteCloser) String() string {
	return bwc.buf.String()
}

func TestAsynchronousWriter_Write(t *testing.T) {
	buf := NewBufferWriteCloser()
	ctx, cancel := context.WithCancel(context.Background())
	writer := NewAsyncWriter(ctx, buf)

	wg := &sync.WaitGroup{}
	writer.Start(wg)

	data := []byte("Hello, World!")

	_, err := writer.Write(data)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond) // Allow some time for the write to complete asynchronously.

	assert.Equal(t, string(data), buf.String())

	cancel() // Stop the writer goroutine.

	wg.Wait() // Wait for the writer goroutine to finish.

	_, err = writer.Write(data)
	assert.EqualError(t, err, ErrClosed.Error(), "Expected ErrClosed after writer is closed")
}

func TestAsynchronousWriter_ErrorChan(t *testing.T) {
	buf := NewBufferWriteCloser()
	ctx, cancel := context.WithCancel(context.Background())
	writer := NewAsyncWriter(ctx, buf)

	wg := &sync.WaitGroup{}
	writer.Start(wg)

	data := []byte("Hello, World!")

	// Simulate an error on the writer goroutine.
	writer.errChan <- errors.New("write error")

	_, err := writer.Write(data)
	assert.EqualError(t, err, "write error", "Expected write error from the error channel")

	cancel() // Stop the writer goroutine.

	wg.Wait() // Wait for the writer goroutine to finish.

	_, err = writer.Write(data)
	assert.EqualError(t, err, ErrClosed.Error(), "Expected ErrClosed after writer is closed")
}

// async_test.go

func TestAsynchronousWriter_QueueFull(t *testing.T) {
	buf := NewDelayedWriter(&bytes.Buffer{}, 2*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	writer := NewAsyncWriter(ctx, buf)

	wg := &sync.WaitGroup{}
	writer.Start(wg)

	// Fill up the queue with messages until it reaches its maximum capacity.
	for i := 0; i <= QueueSize; i++ {
		data := []byte(fmt.Sprintf("Message %d\n", i))
		_, err := writer.Write(data)
		assert.NoError(t, err)
	}

	// The next write should return ErrQueueFull.
	data := []byte("Queue Full!")
	_, err := writer.Write(data)
	assert.EqualError(t, err, ErrQueueFull.Error(), "Expected ErrQueueFull when the queue is full")

	cancel() // Stop the writer goroutine.

	wg.Wait() // Wait for the writer goroutine to finish.
}

func TestAsynchronousWriterWithErrorAfterWriter(t *testing.T) {
	// Create an ErrorAfterWriter that returns an error after 3 write operations.
	buf := &bytes.Buffer{}
	maxWriteCount := 5
	errAfterWrites := 3
	errorWriter := NewErrorAfterWriter(buf, maxWriteCount, errAfterWrites)

	ctx, cancel := context.WithCancel(context.Background())
	writer := NewAsyncWriter(ctx, errorWriter)

	wg := &sync.WaitGroup{}
	writer.Start(wg)

	data := []byte("Hello, Error After Writer!")

	// Write 3 times before the ErrorAfterWriter returns an error.
	for i := 0; i <= errAfterWrites; i++ {
		n, err := writer.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)
	}

	time.Sleep(100 * time.Millisecond)
	// The fourth write should return an error from the ErrorAfterWriter.
	_, err := writer.Write(data)
	assert.EqualError(t, err, "error after write operations", "Expected error after write operations")

	// Wait for the writer goroutine to finish.
	cancel()
	wg.Wait()

	// Attempt to write again after the writer is closed, it should return ErrClosed.
	_, err = writer.Write(data)
	assert.EqualError(t, err, ErrClosed.Error(), "Expected ErrClosed after writer is closed")
}

func TestAsynchronousWriter_OnClose(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	underlyingWriter := NewErrorAfterWriter(buf, 5, 3) // Error after 3 writes
	writer := NewAsyncWriter(ctx, underlyingWriter)

	wg := &sync.WaitGroup{}
	writer.Start(wg)

	data := []byte("Hello, World!")

	// Write 2 times to the writer.
	for i := 0; i < 2; i++ {
		n, err := writer.Write(data)
		assert.NoError(t, err)
		assert.Equal(t, len(data), n)
	}

	// Now, close the writer. This will trigger the onClose() function,
	// and it will flush the remaining buffered data to the underlying writer.
	writer.Close()

	// Wait for the writer goroutine to finish.
	wg.Wait()

	// Ensure the underlying writer received the data.
	assert.Equal(t, string(data)+string(data), buf.String())

	// Attempt to write again after the writer is closed, it should return ErrClosed.
	_, err := writer.Write(data)
	assert.EqualError(t, err, ErrClosed.Error(), "Expected ErrClosed after writer is closed")
	cancel() // to avoid context leak
	// Check if the error from the onClose() function is sent through the errChan.
	// errFromChan := <-writer.errChan
	// assert.Error(t, errFromChan, "Expected error on the errChan from onClose()")
}

func TestAsynchronousWriter_Close(t *testing.T) {
	buf := NewBufferWriteCloser()
	ctx, cancel := context.WithCancel(context.Background())
	underlyingWriter := buf
	writer := NewAsyncWriter(ctx, underlyingWriter)

	wg := &sync.WaitGroup{}
	writer.Start(wg)

	// Close the writer for the first time. It should return nil.
	err := writer.onClose()
	assert.NoError(t, err, "Expected no error on the first Close() call")

	// Attempt to close the writer again. It should return ErrClosed.
	err = writer.onClose()
	assert.EqualError(t, err, ErrClosed.Error(), "Expected ErrClosed on the second Close() call")

	// Cancel the context and wait for the writer goroutine to finish.
	cancel()
	wg.Wait()

	// Attempt to write after closing, it should return ErrClosed.
	data := []byte("Hello, World!")
	_, err = writer.Write(data)
	assert.EqualError(t, err, ErrClosed.Error(), "Expected ErrClosed after writer is closed")
}

type errorWriter struct{}

func (ew *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write error")
}

func (ew *errorWriter) Close() error {
	return errors.New("close error")
}

func TestAsynchronousWriter_OnCloseWithError(t *testing.T) {
	// Create an AsynchronousWriter with an errorWriter that always returns an error when writing.
	ctx, cancel := context.WithCancel(context.Background())
	writer := NewAsyncWriter(ctx, &errorWriter{})
	defer cancel()

	// Write some data to the writer.
	data := []byte("Hello, World!")
	n, err := writer.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	// Cancel the context to trigger the writer to stop and call onClose() to flush remaining data.
	//cancel()

	// Since the writer always returns an error when writing, onClose() should handle the error and return nil.
	err = writer.onClose()
	assert.Error(t, err, "Expected onClose() to handle write error gracefully")
}

type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestAsyncWriter_Close(t *testing.T) {
	mockWriter := new(MockWriter)
	mockWriter.On("Close").Return(nil)

	// Create a context with cancel to control the AsyncWriter.
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the AsyncWriter with the mock Writer and context.
	asyncWriter := NewAsyncWriter(ctx, mockWriter)

	// Create a wait group to wait for the writer goroutine to stop.
	wg := &sync.WaitGroup{}

	// Start the asynchronous writer goroutine.
	asyncWriter.Start(wg)

	// Call the Close() function.
	err := asyncWriter.Close()

	// Wait for the writer goroutine to stop.
	wg.Wait()

	// Verify that Close() returned no error.
	assert.NoError(t, err)

	// Verify that the Close() method of the mock Writer was called.
	mockWriter.AssertCalled(t, "Close")

	cancel()
}

func TestAsyncWriter_Close_Error(t *testing.T) {
	expectedErr := ErrClosed
	mockWriter := new(MockWriter)
	mockWriter.On("Close").Return(expectedErr)

	// Create a context with cancel to control the AsyncWriter.
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the AsyncWriter with the mock Writer and context.
	asyncWriter := NewAsyncWriter(ctx, mockWriter)

	// Create a wait group to wait for the writer goroutine to stop.
	wg := &sync.WaitGroup{}

	// Start the asynchronous writer goroutine.
	asyncWriter.Start(wg)

	cancel() // Stop the writer goroutine.
	//we need to wait to make sure the writer goroutine is finished
	// otherwise the test will fail because Close() will return nil
	time.Sleep(1 * time.Millisecond)
	// Call the Close() function.
	err := asyncWriter.Close()

	// Wait for the writer goroutine to stop.
	wg.Wait()

	// Verify that the Close() method of the mock Writer was called.
	mockWriter.AssertCalled(t, "Close")

	// Verify that Close() returned the expected error.
	assert.EqualError(t, err, expectedErr.Error())

}
