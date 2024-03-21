package ioutils

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// GzipReader implements io.Reader and io.Closer interfaces
// for reading both gziped and non-gziped files.
type GzipReader struct {
	r    io.ReadCloser
	orig io.ReadCloser
}

// NewGzipReader creates a GzipReader from a file path.
func NewGzipReader(filePath string) (*GzipReader, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	gr := &GzipReader{
		orig: file,
		r:    nil,
	}

	if isGzipped(file) {
		gr.r, err = gzip.NewReader(file)
		if err != nil {
			// to keep the original error
			_ = file.Close()
			return nil, err
		}
	}

	return gr, nil
}

// isGzipped checks if a file is gziped.
func isGzipped(file *os.File) bool {
	buf := make([]byte, 2)
	_, err := file.Read(buf)
	if err != nil {
		return false
	}
	_, err = file.Seek(0, io.SeekStart) // Reset the file position
	return err == nil && buf[0] == 0x1f && buf[1] == 0x8b
}

// Read implements the io.Reader interface.
func (gr *GzipReader) Read(p []byte) (n int, err error) {
	if gr.r != nil {
		return gr.r.Read(p)
	}
	return gr.orig.Read(p)
}

// Close implements the io.Closer interface.
func (gr *GzipReader) Close() error {
	if gr.r != nil {
		//#nosec G104 -- We need to close main one in anyway.
		gr.r.Close()
	}
	return gr.orig.Close()
}

// GzipWriter implements io.Writer and io.Closer interfaces
// for writing both gziped and non-gziped files.
type GzipWriter struct {
	w    *gzip.Writer
	orig io.WriteCloser
}

// NewGzipWriter creates a GzipWriter for the specified file path.
// If gzipOutput is true, the output will be gziped; otherwise, it will not be gziped.
func NewGzipWriter(filePath string, gzipOutput bool) (*GzipWriter, error) {
	file, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	gw := &GzipWriter{
		orig: file,
		w:    nil,
	}

	if gzipOutput {
		gw.w = gzip.NewWriter(file)
	}

	return gw, nil
}

// Write implements the io.Writer interface.
func (gw *GzipWriter) Write(p []byte) (n int, err error) {
	if gw.w != nil {
		return gw.w.Write(p)
	}
	return gw.orig.Write(p)
}

// Close implements the io.Closer interface.
func (gw *GzipWriter) Close() error {
	if gw.w != nil {
		// Flush the gzip writer
		//#nosec G104 -- We need to close it anyway.
		gw.w.Flush()
		//#nosec G104 -- We need to close the main file anyway.
		gw.w.Close()
	}
	return gw.orig.Close()
}

// GzipWriteFile writes the given byte slice to a file. If the `compress` flag is true,
// it compresses the file using gzip and appends the ".gz" extension to the file name.
//
// Parameters:
// - path: the path to the file to be written.
// - b: the byte slice to be written to the file.
// - compress: a flag indicating whether the file should be compressed using gzip.
//
// Returns:
// - The number of bytes written to the file.
// - An error if any occurred.
func GzipWriteFile(path string, b []byte, compress bool) (int, error) {
	if compress {
		path = path + ".gz"
	}

	f, err := NewGzipWriter(filepath.Clean(path), compress)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	// Write the structure to an msgpack file
	return f.Write(b)
}
