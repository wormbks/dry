package dry

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"os"
)

type GzipFileReader interface {
	GetReader() (io.Reader, error)
	Close() error
}

type gzipFileReader struct {
	file       *os.File
	reader     io.Reader
	gzipReader *gzip.Reader
}

// NewFileReader returns a new FileReader for a file given its path.
// It opens the file specified by the filePath and returns a FileReader interface for it.
func NewGzipFileReader(filePath string) (GzipFileReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &gzipFileReader{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

// GetReader returns a reader for the file, either normal or gzip-compressed.
//
// Please note that this method relies on the GZIP magic number to detect
// GZIP files. If the file does not have the correct magic number, it will be treated
// as a normal file.
func (f *gzipFileReader) GetReader() (io.Reader, error) {
	if f.gzipReader != nil {
		return f.gzipReader, nil
	}

	peekBytes := make([]byte, 2)
	_, err := f.reader.Read(peekBytes)
	if err != nil {
		if err == io.EOF {
			// Reset the file offset to the beginning
			if _, seekErr := f.file.Seek(0, io.SeekStart); seekErr != nil {
				return nil, seekErr
			}
			return f.file, nil
		}
		return nil, err
	}

	// Check if the file has the GZIP magic number
	if peekBytes[0] == 0x1f && peekBytes[1] == 0x8b {
		gzipReader, err := gzip.NewReader(io.MultiReader(bytes.NewReader(peekBytes), f.reader))
		if err != nil {
			return nil, err
		}
		f.gzipReader = gzipReader
		return f.gzipReader, nil
	}

	// Reset the file offset to the beginning
	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	return f.file, nil
}

// Close closes the file and the gzip reader if applicable.
func (f *gzipFileReader) Close() error {
	var err error
	if f.gzipReader != nil {
		err = f.gzipReader.Close()
		f.gzipReader = nil
	}
	if closeErr := f.file.Close(); err == nil {
		err = closeErr
	}
	return err
}
