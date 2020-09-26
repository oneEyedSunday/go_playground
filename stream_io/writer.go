package main

import (
	"bytes"
	"fmt"
	"os"
)

// DefaultWriterFunc writes into a buffer
func DefaultWriterFunc() {
	proverbs := []string{
		"Channels orchestrate mutexes serialize",
		"Cgo is not Go",
		"Errors are values",
		"Don't panic",
	}

	var writer bytes.Buffer

	for _, p := range proverbs {
		n, err := writer.Write([]byte(p))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if n != len(p) {
			fmt.Println("Failed to write data")
			os.Exit(1)
		}
	}

	fmt.Println(writer.String())
}
