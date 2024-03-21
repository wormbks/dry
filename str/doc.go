package str

/*
Package str provides various utilities for manipulating and processing strings.

It includes a StringBuilder struct, which allows efficient construction of strings.
StringBuilder has methods for appending a variety of types including strings, bytes,
integers, floats, and booleans. It also provides methods for converting the string
builder to a string or a byte slice, resetting the string builder, growing the size
of the internal buffer, and writing to an io.Writer.

This package is intended to provide more efficient string concatenation than using
the '+' operator, and more convenience than using the fmt package for basic string
manipulation tasks.
*/
