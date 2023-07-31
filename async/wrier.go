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
	wr       io.Writer
	ctx      context.Context
	queue    chan *bytes.Buffer
	errChan  chan error
	isClosed atomic.Bool
	stop     chan int // Channel to signal the writer goroutine to stop
}

func NewAsyncWriter(context context.Context, writer io.WriteCloser) *AsyncWriter {
	result := &AsyncWriter{
		wr:      writer,
		ctx:     context,
		queue:   make(chan *bytes.Buffer, QueueSize),
		errChan: make(chan error, QueueSize),
		stop:    make(chan int),
	}
	result.isClosed.Store(false)
	// buffer pool for asynchronous writer
	return result
}

// start the asynchronous writer
func (w *AsyncWriter) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go w.writer(wg)
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
func (w *AsyncWriter) writer(wg *sync.WaitGroup) {
	var err error
	defer wg.Done()
	for {
		select {
		case b := <-w.queue:
			// fmt.Printf("write %s \n", string(b.Bytes()))
			_, err = w.wr.Write(b.Bytes())
			w.sendIfError(err)
			_asyncBufferPool.Put(b)
		case <-w.ctx.Done():
			// Stop the writer goroutine gracefully when the context is canceled.
			w.onClose()
			return
		case <-w.stop:
			// Stop the writer goroutine gracefully when the stop signal is received.
			w.onClose()
			return
		}
	}
}

func (w *AsyncWriter) sendIfError(err error) {
	if err != nil {
		select {
		case w.errChan <- err:
		default:
		}
	}
}

func (w *AsyncWriter) Close() (err error) {
	if !w.isClosed.Load() {
		//close(w.stop) // Send the stop signal to the writer goroutine
		w.stop <- 1
		return err
	}
	return ErrClosed
}

// onClose set closed and close the file once
func (w *AsyncWriter) onClose() (err error) {
	if w.isClosed.Load() {
		return ErrClosed
	}
	w.isClosed.Store(true)
	w.flushQueue()
	// does writer has ioCloser interface
	if w, ok := w.wr.(io.Closer); ok {
		err = w.Close()
	}
	return err
}

// flushQueue process remaining buffered data for asynchronous writer
func (w *AsyncWriter) flushQueue() {
	var err error
	for {
		select {
		case b := <-w.queue:
			// flush all remaining field
			if _, err = w.wr.Write(b.Bytes()); err != nil {
				select {
				case w.errChan <- err:
				default:
				}
			}
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
