package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"oneeyedsunday/streamio/readers"
	"os"
	"strings"
)

// ReaderConfig is configuration for my reader
type ReaderConfig struct {
	successMsg string
	splitFunc  bufio.SplitFunc
}

// ChannelWriter writes bytes to a channel
type ChannelWriter struct {
	ch chan byte
}

func streamFromString(source string, useDigitReader bool) {
	var reader io.Reader = strings.NewReader(source)

	if useDigitReader {
		reader = readers.NewDigitReader(reader)
	}

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

func streamFromFile(path string, asBinary bool) {
	file, err := os.Open(path)
	if err != nil {
		handleFileOpenError(err)
		return
	}

	defer file.Close()

	successMsg := "Successfully opened %s"
	var buffer []byte = make([]byte, 20)

	if asBinary {
		successMsg += "; reading as Binary"
		buffer = make([]byte, 256)
	}

	successMsg += "\n"

	fmt.Printf(successMsg, file.Name())

	reader := bufio.NewReader(file)
	for {
		numRead, err := reader.Read(buffer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("Fatal error occured: %s\n", err)
			}
			break
		}

		if asBinary {
			fmt.Printf("[+]%s", hex.Dump(buffer[:numRead]))
		} else {
			fmt.Printf("%s\n-----------------------------------------\n", string(buffer[:numRead]))
		}
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

// customWriterFunc runs a writer that writes to a channel
func customWriterFunc() {
	writer := NewCustomWriter()
	go func() {
		defer writer.Close()
		writer.Write([]byte("Jaja "))
		writer.Write([]byte("Moi"))
	}()

	for c := range writer.Chan() {
		fmt.Printf("%c", c)
	}

	fmt.Println()
}

func main() {
	sourceString := flag.String("string_source", "", "Stream from string")
	useDigitReader := flag.Bool("digit", false, "Use Custom Digit Reader")
	filePath := flag.String("file_source", "", "Stream from file source")
	asBinary := flag.Bool("bin", false, "Read as Binary")
	lineByLine := flag.Bool("line", false, "Read line by line")
	wordByWord := flag.Bool("word", false, "Read word by word")
	defaultWriter := flag.Bool("defaultWriter", false, "Demo default writer")
	chanWriter := flag.Bool("chanWriter", false, "Demo Channel Writer")
	fileIO := flag.Bool("files", false, "Demo Files IO")
	stdIO := flag.Bool("std", false, "Std IO")
	flag.Parse()

	proverbs := []string{
		"Channels orchestrate mutexes serialize\n",
		"Cgo is not Go\n",
		"Errors are values\n",
		"Don't panic\n",
	}

	switch {
	case *sourceString != "":
		streamFromString(*sourceString, *useDigitReader)
	case *filePath != "" && *wordByWord:
		streamFromFileWithSplit(*filePath, ReaderConfig{successMsg: "Successfully opened %s to read word by word\n", splitFunc: bufio.ScanWords})
	case *filePath != "" && *lineByLine:
		streamFromFileWithSplit(*filePath, ReaderConfig{successMsg: "Successfully opened %s to read line by line\n", splitFunc: bufio.ScanLines})
	case *filePath != "":
		streamFromFile(*filePath, *asBinary)
	case *defaultWriter:
		DefaultWriterFunc()
	case *chanWriter:
		customWriterFunc()
	case *fileIO:
		{
			file, err := os.Create("./proverbs.txt")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()

			for _, p := range proverbs {
				n, err := file.Write([]byte (p))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if n != len(p) {
					fmt.Println("failed to write data")
					os.Exit(1)
				}
			}

			fmt.Println("file write done")
		}
	case *stdIO:
		{
			for _, p := range proverbs {
				n, err := os.Stdout.Write([]byte(p))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if n != len(p) {
					fmt.Println("failed to write data")
					os.Exit(1)
				}
			}
		}
	default:
		fmt.Printf("No valid selection made. Run %s -help\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Println("Walthrough over.")
}
