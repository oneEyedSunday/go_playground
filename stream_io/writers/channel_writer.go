package writers

// ChannelWriter writes bytes to a channel
type ChannelWriter struct {
	ch chan byte
}

// NewCustomWriter does yada yada
func NewCustomWriter() *ChannelWriter {
	return &ChannelWriter{make(chan byte, 1024)}
}

// Chan exposes the byte channel
func (w *ChannelWriter) Chan() <-chan byte {
	return w.ch
}

// Write implements the Writer interface
func (w *ChannelWriter) Write(p []byte) (int, error) {
	n := 0
	for _, b := range p {
		w.ch <- b
		n++
	}

	return n, nil
}

// Close closes
func (w *ChannelWriter) Close() error {
	close(w.ch)
	return nil
}
