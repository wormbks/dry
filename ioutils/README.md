# IO Utilities Package

- [IO Utilities Package](#io-utilities-package)
  - [Gzip Reader](#gzip-reader)
    - [`GzipReader`](#gzipreader)
    - [Methods](#methods)
  - [Gzip Writer](#gzip-writer)
    - [`GzipWriter`](#gzipwriter)
    - [Methods](#methods-1)
  - [File Writing Utilities](#file-writing-utilities)
    - [`GzipWriteFile`](#gzipwritefile)
  - [Example Usage](#example-usage)
    - [Notes](#notes)
    - [References](#references)
  - [Folder Structure Creation](#folder-structure-creation)
    - [`CreateFolderStructure`](#createfolderstructure)
  - [Example Usage](#example-usage-1)
    - [Notes](#notes-1)
    - [References](#references-1)


The `ioutils` package offers utilities for reading from and writing to files, with support for both gzip-compressed and uncompressed files. This package includes types and functions to simplify file I/O operations in Go.

## Gzip Reader

### `GzipReader`

The `GzipReader` struct implements the `io.Reader` and `io.Closer` interfaces for reading both gzip-compressed and uncompressed files.

- **NewGzipReader(filePath string) (*GzipReader, error)**: Creates a `GzipReader` from a file path.

### Methods

- **Read(p []byte) (n int, err error)**: Implements the `io.Reader` interface for reading data from the file.
- **Close() error**: Implements the `io.Closer` interface for closing the file.

## Gzip Writer

### `GzipWriter`

The `GzipWriter` struct implements the `io.Writer` and `io.Closer` interfaces for writing both gzip-compressed and uncompressed files.

- **NewGzipWriter(filePath string, gzipOutput bool) (*GzipWriter, error)**: Creates a `GzipWriter` for the specified file path.

### Methods

- **Write(p []byte) (n int, err error)**: Implements the `io.Writer` interface for writing data to the file.
- **Close() error**: Implements the `io.Closer` interface for closing the file.

## File Writing Utilities

### `GzipWriteFile`

- **GzipWriteFile(path string, b []byte, compress bool) (int, error)**: Writes a byte slice to a file, optionally compressing it using gzip and appending the ".gz" extension to the file name if `compress` is `true`.

## Example Usage

```go
// Read from a gzip-compressed file
reader, err := ioutils.NewGzipReader("example.gz")
if err != nil {
    // Handle error
}
defer reader.Close()

// Write to a file, compressing it if needed
data := []byte("Hello, world!")
_, err = ioutils.GzipWriteFile("output.txt", data, true)
if err != nil {
    // Handle error
}
```
### Notes

- The package provides convenience functions for working with gzip-compressed files without the need for manual handling of compression and decompression.
- Ensure proper error handling and resource cleanup when using file I/O operations to prevent resource leaks.

### References

- [Go Standard Library - compress/gzip](https://golang.org/pkg/compress/gzip/)
- [Go Standard Library - os](https://golang.org/pkg/os/)
- [Go Standard Library - path/filepath](https://golang.org/pkg/path/filepath/)

For more detailed usage examples and documentation, refer to the package source code and Go documentation.


## Folder Structure Creation

### `CreateFolderStructure`

The `CreateFolderStructure` function creates a folder structure based on the provided path. It creates all necessary parent directories if they don't already exist.

- **CreateFolderStructure(path string) error**: Creates a folder structure at the specified path.

## Example Usage

```go
// Create a folder structure
err := ioutils.CreateFolderStructure("/path/to/new/folder")
if err != nil {
    // Handle error
}
```

### Notes

- The function uses `os.MkdirAll` to create directories recursively, ensuring that all parent directories are created if they don't exist.
- Ensure proper error handling to handle cases where directory creation fails.

### References

- [Go Standard Library - os](https://golang.org/pkg/os/)

For more detailed usage examples and documentation, refer to the package source code and Go documentation.


