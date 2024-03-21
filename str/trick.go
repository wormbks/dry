package str

// It is based on "140x Faster String to Byte and Byte to String Conversions with Zero Allocation in Go"
// https://josestg.medium.com/140x-faster-string-to-byte-and-byte-to-string-conversions-with-zero-allocation-in-go-200b4d7105fc

import (
	"bytes"
	"errors"
	"unsafe"
)

// BytesToString converts bytes to a string without memory allocation.
// NOTE: The given bytes MUST NOT be modified since they share the same backing array
// with the returned string.
func BytesToString(b []byte) string {
	// Ignore if your IDE shows an error here; it's a false positive.
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}

// StringToBytes converts a string to a byte slice without memory allocation.
// NOTE: The returned byte slice MUST NOT be modified since it shares the same backing array
// with the given string.
func StringToBytes(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}

// RemoveTabsAndNewlines optimizes removing tab and newline characters from a byte slice.
func RemoveTabsAndNewlines(src []byte, buf *bytes.Buffer) error {
	if src == nil {
		return errors.New("source slice is nil")
	}

	for i := 0; i < len(src); i++ {
		b := src[i]
		if b != 9 && b != 10 && b != 13 {
			buf.WriteByte(b)
		}
	}

	return nil
}
