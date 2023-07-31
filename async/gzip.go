package async

import (
	"compress/gzip"
	"context"
	"os"
	"sync"
)

// GzipFileWriter represents a file writer that compresses data using gzip and writes it asynchronously.
type GzipFileWriter struct {
	writer     *AsyncWriter // The underlying AsyncWriter.
	gzipWriter *gzip.Writer // Gzip writer to compress data.
	file       *os.File
	wg         sync.WaitGroup
}

// NewGzipFileWriter creates a new AsyncGzipFileWriter with the specified file path.
func NewGzipFileWriter(filePath string, gzipIt bool) (*GzipFileWriter, error) {
	ctx := context.Background()
	// Open the file for writing.
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	res := &GzipFileWriter{
		file: file,
	}

	if gzipIt {
		// Create a gzip writer that wraps the file.
		res.gzipWriter = gzip.NewWriter(file)
		// Create an AsynchronousWriter that wraps the gzip writer.
		res.writer = NewAsyncWriter(ctx, res.gzipWriter, &res.wg)
	} else {
		// Create an AsynchronousWriter that wraps the file.
		res.writer = NewAsyncWriter(ctx, file, &res.wg)
	}

	return res, err
}

// Write writes the compressed data asynchronously.
func (aw *GzipFileWriter) Write(data []byte) (int, error) {
	return aw.writer.Write(data)
}

// Close closes the gzip writer and flushes any remaining buffered data.
func (aw *GzipFileWriter) Close() (err error) {
	// Close the underlying AsynchronousWriter.
	err = aw.writer.Close(&aw.wg)
	// // If  use gzip writer, close it.
	// // Otherwise, it just was closed by AsyncWriter.
	if aw.gzipWriter != nil {
		err = aw.gzipWriter.Flush()
		if err != nil {
			return err
		}
		err = aw.gzipWriter.Close()
		if err != nil {
			return err
		}
	}
	if aw.file != nil {
		err = aw.file.Close()
	}

	return err
}
