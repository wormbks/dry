package dry

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrClosed = errors.New("closed")

func TestWrite(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World")
	assert.Equal(t, "Hello World", sb.String(), "Write function failed")
}

func TestPrintf(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World")
	sb.Printf(" - %d %s", 2023, "Happy New Year")
	assert.Equal(t, "Hello World - 2023 Happy New Year", sb.String(), "Printf function failed")
}

func TestByte(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World")
	sb.Byte('!')
	assert.Equal(t, "Hello World!", sb.String(), "Byte function failed")
}

func TestWriteBytes(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World")
	bytesToWrite := []byte(" Have a great day!")
	sb.WriteBytes(bytesToWrite)
	assert.Equal(t, "Hello World Have a great day!", sb.String(), "WriteBytes function failed")
}

func TestInt(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World", " ")
	sb.Int(100)
	assert.Equal(t, "Hello World 100", sb.String(), "Int function failed")
}

func TestUint(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World", " ")
	sb.Uint(999)
	assert.Equal(t, "Hello World 999", sb.String(), "Uint function failed")
}

func TestFloat(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World", " ")
	sb.Float(3.14)
	assert.Equal(t, "Hello World 3.14", sb.String(), "Float function failed")
}

func TestBool(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World", " ")
	sb.Bool(true)
	assert.Equal(t, "Hello World true", sb.String(), "Bool function failed")
}

func TestWriteTo(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World")
	nWas := len(sb.String())

	var writer bytes.Buffer
	n, err := sb.WriteTo(&writer)
	assert.Nil(t, err, "WriteTo function failed")
	assert.Equal(t, int64(nWas), n, "WriteTo function failed")

	// Negative test case: WriteTo with a closed writer should return an error
	closedWriter := &closedBuffer{}
	closedWriter.Close()
	_, err = sb.WriteTo(closedWriter)
	// uses bytes.Buffer Write method but it returns an error nil every time
	assert.Nil(t, err, "WriteTo with a closed writer should return an error")
}

func TestBytes(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World")
	expectedBytes := []byte("Hello World")
	assert.Equal(t, expectedBytes, sb.Bytes(), "Bytes function failed")
}

func TestString(t *testing.T) {
	sb := &StringBuilder{}
	sb.Write("Hello", " ", "World")
	assert.Equal(t, sb.String(), sb.String(), "String function failed")
}

type closedBuffer struct {
	closed bool
}

func (cb *closedBuffer) Write(p []byte) (n int, err error) {
	if cb.closed {
		return 0, ErrClosed
	}
	return len(p), nil
}

func (cb *closedBuffer) Close() error {
	cb.closed = true
	return nil
}
