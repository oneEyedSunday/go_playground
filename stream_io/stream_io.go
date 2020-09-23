package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

func streamFromString(source string) {
	reader := strings.NewReader(source)
	buffer := make([]byte, 10)
	for {
		numRead, err := reader.Read(buffer)
		if errors.Is(err, io.EOF) {
			fmt.Printf("End of file. Outputting any remaining bytes... %s\n", string(buffer[:numRead]))
			break
		} else if err != nil {
			// non null err that isnt EOF
			fmt.Errorf("An unknown error occured: %s", err)
		}

		fmt.Printf("Read %d chars which were: %s\n", numRead, string(buffer[:numRead]))
	}
}

func main() {
	sourceString := flag.String("string_source", "", "Run stream from string demo")
	flag.Parse()

	switch {
	case *sourceString != "":
		streamFromString(*sourceString)
		// fallthrough
	default:
		fmt.Println("No valid selection made.")
	}

	fmt.Println("Walthrough over.")
}
