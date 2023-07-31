package async

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

var (
	// BufferSize defined the buffer size, by default 1 KB buffer will be allocated
	BufferSize = 1024
	// QueueSize defined the queue size for asynchronous write
	QueueSize = 1024
	// Precision defined the precision about the reopen operation condition
	// check duration within second
	Precision = 1
	// DefaultFileMode set the default open mode rw-r--r-- by default
	DefaultFileMode = os.FileMode(0o644)
	// DefaultFileFlag set the default file flag
	DefaultFileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND

	// ErrInternal defined the internal error
	ErrInternal = errors.New("error internal")
	// ErrClosed defined write while ctx close
	ErrClosed = errors.New("error write on close")
	// ErrInvalidArgument defined the invalid argument
	ErrInvalidArgument = errors.New("error argument invalid")
	// ErrQueueFull defined the queue full
	ErrQueueFull = errors.New("async log queue full")

	ErrWrongType = errors.New("wrong type")
)

type AsyncWriter struct {
	wr         io.Writer
	ctx        context.Context
	queue      chan *bytes.Buffer
	errChan    chan error
	isClosed   atomic.Bool
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
}

func NewAsyncWriter(ctx context.Context, writer io.WriteCloser, wg *sync.WaitGroup) *AsyncWriter {
	result := &AsyncWriter{
		wr:      writer,
		queue:   make(chan *bytes.Buffer, QueueSize),
		errChan: make(chan error, QueueSize),
		wg:      wg,
		// stop:    make(chan int),
	}
	result.ctx, result.cancelFunc = context.WithCancel(ctx)
	result.isClosed.Store(false)
	// buffer pool for asynchronous writer
	result.start()
	return result
}

// start the asynchronous writer
func (w *AsyncWriter) start() {
	w.wg.Add(1)
	go w.writer()
}

// Only when the error channel is empty, otherwise nothing will write and the last error will be
// returned the error channel
func (w *AsyncWriter) Write(b []byte) (int, error) {
	if !w.isClosed.Load() {
		ok := false
		for !ok {
			select {
			case err := <-w.errChan:
				// NOTE this error caused by last write maybe ignored
				return 0, err
			default:
				ok = true
			}
		}

		bb := _asyncBufferPool.Get().(*bytes.Buffer)
		bb.Reset()          // remove old buffer data
		n, _ := bb.Write(b) // bytes.Buffer Write returns error nil	all the time
		select {

		case w.queue <- bb:
			return n, nil
		default:
			return 0, ErrQueueFull
		}
	}

	return 0, ErrClosed
}

// writer do the asynchronous write independently
// Take care of reopen, I am not sure if there need no lock
func (w *AsyncWriter) writer() {
	var err error
	defer w.wg.Done()
	for {
		select {
		case b := <-w.queue:
			_, err = w.wr.Write(b.Bytes())
			w.sendIfError(err)
			_asyncBufferPool.Put(b)
		case <-w.ctx.Done():
			// Stop the writer goroutine gracefully when the context is canceled.
			w.onClose()
			return
		}
	}
}

// sendIfError sends the error to the error channel if it is not nil.
//
// It takes an error as a parameter and does not return anything.
func (w *AsyncWriter) sendIfError(err error) {
	if err != nil {
		select {
		case w.errChan <- err:
		default:
		}
	}
}

// Close closes the AsyncWriter.
//
// It cancels the writer goroutine and waits for it to finish.
// It takes a sync.WaitGroup as a parameter to coordinate the closing.
// It returns an error if there is any.
func (w *AsyncWriter) Close(wg *sync.WaitGroup) (err error) {
	// close(w.stop) // Send the stop signal to the writer goroutine
	w.cancelFunc()
	wg.Wait()
	return nil
}

// onClose set closed and close the file once
func (w *AsyncWriter) onClose() (err error) {
	if w.isClosed.Load() {
		return ErrClosed
	}
	w.isClosed.Store(true)
	w.flushQueue()
	// does underlining writer has io.Closer interface
	// if w, ok := w.wr.(io.Closer); ok {
	// 	err = w.Close()
	// }
	return err
}

// flushQueue process remaining buffered data for asynchronous writer
func (w *AsyncWriter) flushQueue() {
	var err error
	for {
		select {
		case b := <-w.queue:
			// flush all remaining field
			_, err = w.wr.Write(b.Bytes())
			w.sendIfError(err)
			_asyncBufferPool.Put(b)
		default: // after the queue was empty, return
			return
		}
	}
}

var _asyncBufferPool = sync.Pool{
	New: func() interface{} {
		// return make([]byte, BufferSize)
		return bytes.NewBuffer(make([]byte, 0, BufferSize))
	},
}
