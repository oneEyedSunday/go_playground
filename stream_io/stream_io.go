package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReaderConfig is configuration for my reader
type ReaderConfig struct {
	successMsg string
	splitFunc  bufio.SplitFunc
}

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

func handleFileOpenError(err error) {
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
	fmt.Println(fmt.Errorf("A non os error occured: %w\n", err))
}

func streamFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		handleFileOpenError(err)
		return
	}

	defer file.Close()

	fmt.Printf("Successfully opened %s\n", file.Name())

	buffer := make([]byte, 20)
	reader := bufio.NewReader(file)
	for {
		numRead, err := reader.Read(buffer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("Fatal error occured: %s\n", err)
			}
			break
		}

		fmt.Printf("%s\n-----------------------------------------\n", string(buffer[:numRead]))
	}
	fmt.Println()

}

func streamFromFileWithSplit(path string, cfg ReaderConfig) {
	file, err := os.Open(path)

	if err != nil {
		handleFileOpenError(err)
		return
	}

	defer file.Close()

	fmt.Printf(cfg.successMsg, path)

	scanner := bufio.NewScanner(file)
	if cfg.splitFunc != nil {
		scanner.Split(cfg.splitFunc)
	}

	for scanner.Scan() {
		fmt.Printf("[+]\t%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("[-] Error reading: ", err)
	}
}

func main() {
	sourceString := flag.String("string_source", "", "Stream from string")
	filePath := flag.String("file_source", "", "Stream from file source")
	lineByLine := flag.Bool("line", false, "Read line by line")
	wordByWord := flag.Bool("word", false, "Read word by word")
	flag.Parse()

	switch {
	case *sourceString != "":
		streamFromString(*sourceString)
		// fallthrough
	case *filePath != "" && *wordByWord:
		streamFromFileWithSplit(*filePath, ReaderConfig{successMsg: "Successfully opened %s to read word by word\n", splitFunc: bufio.ScanWords})
	case *filePath != "" && *lineByLine:
		streamFromFileWithSplit(*filePath, ReaderConfig{successMsg: "Successfully opened %s to read line by line\n", splitFunc: bufio.ScanLines})
	case *filePath != "":
		streamFromFile(*filePath)
	default:
		fmt.Println("No valid selection made.")
	}

	fmt.Println("Walthrough over.")
}
