package async

import (
	"compress/gzip"
	"context"
	"os"
	"sync"
)

// AsyncGzipFileWriter represents a file writer that compresses data using gzip and writes it asynchronously.
type AsyncGzipFileWriter struct {
	writer     *AsyncWriter // The underlying AsyncWriter.
	gzipWriter *gzip.Writer // Gzip writer to compress data.
	file       *os.File
}

// NewAsyncGzipFileWriter creates a new AsyncGzipFileWriter with the specified file path.
func NewAsyncGzipFileWriter(ctx context.Context, filePath string, gzipOrNot bool) (*AsyncGzipFileWriter, error) {
	// Open the file for writing.
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	var gzipWriter *gzip.Writer
	var asyncWriter *AsyncWriter
	if gzipOrNot {
		// Create a gzip writer that wraps the file.
		gzipWriter = gzip.NewWriter(file)
		// Create an AsynchronousWriter that wraps the gzip writer.
		asyncWriter = NewAsyncWriter(ctx, gzipWriter)
	} else {
		// Create an AsynchronousWriter that wraps the file.
		asyncWriter = NewAsyncWriter(ctx, file)
	}

	return &AsyncGzipFileWriter{
		writer:     asyncWriter,
		gzipWriter: gzipWriter,
		file:       file,
	}, nil
}

func (aw *AsyncGzipFileWriter) Start(wg *sync.WaitGroup) {
	go aw.writer.Start(wg)
}

// Write writes the compressed data asynchronously.
func (aw *AsyncGzipFileWriter) Write(data []byte) (int, error) {
	return aw.writer.Write(data)
}

// Close closes the gzip writer and flushes any remaining buffered data.
func (aw *AsyncGzipFileWriter) Close() (err error) {
	// Close the underlying AsynchronousWriter.
	err = aw.writer.Close()
	// If we don't use gzip writer, close the file.
	// Otherwise, it just was closed by AsyncWriter.
	if err == nil && aw.gzipWriter != nil {
		err = aw.file.Close()
	}

	return err
}
