package dry

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type StringBuilder struct {
	buffer bytes.Buffer
}

// Write concatenates the given strings and appends them to the StringBuilder's buffer.
//
// It takes a variadic parameter of type string.
// It returns a pointer to the StringBuilder.
func (s *StringBuilder) Write(strings ...string) *StringBuilder {
	for _, str := range strings {
		s.buffer.WriteString(str)
	}
	return s
}

// Printf formats and writes the string to the underlying buffer using the
// specified format and arguments. It returns the modified StringBuilder.
//
// format: the format string.
// args: the arguments to be formatted.
// returns: the modified StringBuilder.
func (s *StringBuilder) Printf(format string, args ...interface{}) *StringBuilder {
	fmt.Fprintf(&s.buffer, format, args...)
	return s
}

// Byte appends a byte to the string builder.
//
// value: the byte to be appended.
// returns: the string builder itself.
func (s *StringBuilder) Byte(value byte) *StringBuilder {
	s.buffer.WriteByte(value)
	return s
}

// WriteBytes writes the given bytes to the StringBuilder.
//
// bytes: The bytes to be written.
// Returns: The updated StringBuilder.
func (s *StringBuilder) WriteBytes(bytes []byte) *StringBuilder {
	s.buffer.Write(bytes)
	return s
}

// Int appends an integer value to the string builder.
//
// value: the integer value to be appended.
// Returns a pointer to the string builder.
func (s *StringBuilder) Int(value int) *StringBuilder {
	s.buffer.WriteString(strconv.Itoa(value))
	return s
}

// Uint appends the string representation of the unsigned integer value to the StringBuilder.
//
// value: The unsigned integer value to be appended.
// Returns: The StringBuilder object.
func (s *StringBuilder) Uint(value uint) *StringBuilder {
	s.buffer.WriteString(strconv.FormatUint(uint64(value), 10))
	return s
}

// Float appends a float64 value to the StringBuilder.
//
// value: The float64 value to be appended.
// Returns: A pointer to the StringBuilder.
func (s *StringBuilder) Float(value float64) *StringBuilder {
	s.buffer.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
	return s
}

// Bool appends a boolean value to the string builder.
//
// value: The boolean value to be appended.
// Returns: The updated string builder.
func (s *StringBuilder) Bool(value bool) *StringBuilder {
	s.buffer.WriteString(strconv.FormatBool(value))
	return s
}

// WriteTo writes the contents of the StringBuilder to the given io.Writer.
//
// It returns the number of bytes written and any error encountered.
func (s *StringBuilder) WriteTo(writer io.Writer) (n int64, err error) {
	return s.buffer.WriteTo(writer)
}

// Bytes returns the byte slice representation of the StringBuilder.
//
// No parameters.
// Returns a byte slice.
func (s *StringBuilder) Bytes() []byte {
	return s.buffer.Bytes()
}

// String returns the string representation of the StringBuilder.
//
// No parameters.
// Returns a string.
func (s *StringBuilder) String() string {
	return s.buffer.String()
}
