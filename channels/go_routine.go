package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter = 0
	lock    sync.Mutex
)

func main() {
	fmt.Println("Start")
	go func() {
		fmt.Println("Processing")
	}()

	time.Sleep(time.Millisecond * 10)
	fmt.Println("Done")

	fmt.Println("Without waiting")
	// main process exits without actually running this.
	// needs synchronization
	go func() {
		fmt.Println("Processing")
	}()
	fmt.Println("Done again")

	fmt.Printf("\tWriting to shared var across goroutines. Bad...\n")
	for i := 0; i < 20; i++ {
		go incr()
	}

	time.Sleep(time.Millisecond * 10)

	fmt.Printf("\tUsing locks (mutex) :)\n")

}

func incr() {
	lock.Lock()
	defer lock.Unlock()
	counter++
	fmt.Println(counter)
}
