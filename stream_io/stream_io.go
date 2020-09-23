package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
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
			fmt.Println(fmt.Errorf("An unknown error occured: %w", err))
		}

		fmt.Printf("Read %d chars which were: %s\n", numRead, string(buffer[:numRead]))
	}
}

func streamFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		/*if e, isPError := err.(*os.PathError); isPError {
			fmt.Printf("Failed to open %v: %v\n", e.Path, e.Err)
			// fmt.Println(fmt.Errorf("Error opening path: %v %v\n", path, err))
			return
		}*/
		var pErr *os.PathError
		if errors.As(err, &pErr) {
			fmt.Println(fmt.Errorf("Error opening %v %v\n", pErr.Path, pErr.Err))
			return
		}

		fmt.Println(fmt.Errorf("Another Error %w", err))
		return
	}

	defer file.Close()

	fmt.Printf("Successfully opened %s\n", file.Name())

	buffer := make([]byte, 20)
	reader := io.Reader(file)
	numRead, err := reader(buffer)
	for {
		if err != nil {
			fmt.Errorf("Error: %w\n", err)
			break
		}

		fmt.Printf("[+] %s", string(buffer[:numRead]))
	}

}

func main() {
	sourceString := flag.String("string_source", "", "Stream from string")
	filePath := flag.String("file_source", "", "Stream from file source")
	flag.Parse()

	switch {
	case *sourceString != "":
		streamFromString(*sourceString)
		// fallthrough
	case *filePath != "":
		streamFromFile(*filePath)
	default:
		fmt.Println("No valid selection made.")
	}

	fmt.Println("Walthrough over.")
}
