package readers

import (
	"io"
)

// DigitReader
type DigitReader struct {
	src string
	cur int
}

// NewDigitReader returns a pointer to a DigitReader
func NewDigitReader(src string) *DigitReader {
	return &DigitReader{src: src}
}

func isDigit(r byte) bool {
	return (r >= '0' && r <= '9')
}

func (r *DigitReader) Read(p []byte) (int, error) {
	if r.cur >= len(r.src) {
		return 0, io.EOF
	}

	availToRead := len(r.src) - r.cur
	numRead, bound := 0, 0
	if availToRead >= len(p) {
		// buffer cannot hold all data, only read up to length of buffer into buffer
		bound = len(p)
	} else if availToRead <= len(p) {
		// buffer can hold everything we have left to read
		bound = availToRead
	}

	buffer := make([]byte, bound)
	for numRead < bound {
		if digit := r.src[r.cur]; isDigit(digit) {
			buffer[numRead] = digit
		}
		numRead++
		r.cur++
	}

	copy(p, buffer)
	return numRead, nil

}
