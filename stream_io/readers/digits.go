package readers

import (
	"io"
)

// DigitReaderRaw accepts a string as its source
type DigitReaderRaw struct {
	src string
	cur int
}

// DigitReader accepts an io.Reader
type DigitReader struct {
	reader io.Reader
}

// NewDigitReaderFromSource returns a pointer to a DigitReader
func NewDigitReaderFromSource(src string) *DigitReaderRaw {
	return &DigitReaderRaw{src: src}
}

// NewDigitReader returns a pointer to a DigitReader wrapping an io.Reader
func NewDigitReader(reader io.Reader) *DigitReader {
	return &DigitReader{reader: reader}
}

func isDigit(r byte) bool {
	return (r >= '0' && r <= '9')
}

func (r *DigitReader) Read(p []byte) (int, error) {
	numRead, err := r.reader.Read(p)
	if err != nil {
		return numRead, err
	}
	buffer := make([]byte, numRead)
	for i := 0; i < numRead; i++ {
		if digit := p[i]; isDigit(digit) {
			buffer[i] = digit
		}
	}


	copy(p, buffer)
	return numRead, nil

}
