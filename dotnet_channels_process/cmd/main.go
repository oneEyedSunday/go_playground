package main

import (
	"fmt"

	"github.com/oneeyedsunday/go_playground/dotnet_channels_process/pkg/data"
)

func main() {
	// wire this up to look like the background stuff
	// eventually use a cmd/background run

	c := make(chan data.DataWithKey)

	fmt.Println(len(c))

	// dataProcessor := data.NewDataProcessor()
	// should be an error here
	dataProcessor := new(data.BackgroundDataProcessor)

	dataProcessor.ReadAndFn(func(i data.DataWithKey) error {
		fmt.Printf("Gotten data %v \n", i)
		return nil
	})

	fmt.Println("Here")
}
