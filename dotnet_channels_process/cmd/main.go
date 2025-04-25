package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/oneeyedsunday/go_playground/dotnet_channels_process/pkg/data"
)

func measureExecutionTime(start time.Time) {
	fmt.Printf("After program end, elapsed in %v\n", time.Since(start))
}

func pumpData(ctx context.Context, dataProcessor *data.BackgroundDataProcessor) {
	threeSecondsInterval := time.NewTicker(time.Second * time.Duration(3))
	twoSecondsInterval := time.NewTicker(time.Second * time.Duration(2))

	// push data immediately
	// because Push will block till its read
	// You need to go this, else Execute wont be called
	go dataProcessor.Push(data.DataWithKey{Key: "base"})
	go dataProcessor.Push(data.DataWithKey{Key: "___"})

	go func() {
		// push data after 400 millisecond sor abort if cancelled before that
		select {
		case <-ctx.Done():
			return
		case <-time.Tick(time.Millisecond * time.Duration(400)):
			fmt.Println("attempt to push")
			go dataProcessor.Push(data.DataWithKey{})
			go dataProcessor.Push(data.DataWithKey{Key: "foo", Data: "two"})
			fmt.Println("After pushing")
		}
	}()

	go func() {
		// push data after 5 seconds or abort if cancelled before that
		select {
		case <-ctx.Done():
			return

		case <-time.Tick(time.Second * time.Duration(5)):
			fmt.Println("attempt to push")
			go dataProcessor.Push(data.DataWithKey{Key: "bar", Data: "three"})
			go dataProcessor.Push(data.DataWithKey{Key: "foo", Data: "four"})
			go dataProcessor.Push(data.DataWithKey{})
			fmt.Println("After pushing")
		}
	}()

	go func() {
		defer threeSecondsInterval.Stop()
		defer twoSecondsInterval.Stop()
		<-ctx.Done()
	}()

	func() {
		// push data indefinitely AND abort when cancelled
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopped pumping data indefinitely")
				return
			case <-threeSecondsInterval.C:
				go dataProcessor.Push(data.DataWithKey{Key: "alice", Data: "wonderland"})
			case <-twoSecondsInterval.C:
				go dataProcessor.Push(data.DataWithKey{Key: "keith", Data: "peele"})
			}
		}
	}()

}

func main() {
	startT := time.Now()
	defer measureExecutionTime(startT)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(6))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fmt.Println("Start")

	dataProcessor := data.NewDataProcessor()
	// should be an error here
	// Yes it blocks indefinitely on reading / writing to / from nil channels
	// dataProcessor := new(data.BackgroundDataProcessor)

	go func() {
		fmt.Printf("interrupt %v, cancelling\n", <-interrupt)
		cancel()
	}()

	go pumpData(ctx, dataProcessor)

	// display delayed consumption (and separation of consumption from production)
	<-time.Tick(time.Millisecond * time.Duration(200))

	go dataProcessor.Execute(ctx)

	<-ctx.Done()
}
