package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Worker struct {
	id int
}

func (w *Worker) process(c chan int) {
	for {
		data := <-c
		fmt.Printf("Worker %d got %d\n", w.id, data)
		time.Sleep(time.Millisecond * 1000) // simulate backpressure
	}
}

func main() {
	c := make(chan int, 10)
	// sending c <- 5
	// receiving x := <- c
	// are both blocing

	for i := 0; i < 4; i++ {
		go (&Worker{id: i}).process(c)
	}

	for {
		c <- rand.Int()
		fmt.Printf("Channel length (Buffered) %d capacity is: %d\n", len(c), cap(c))
		time.Sleep(time.Millisecond * 50) // pause exec so user can more easily see flow
	}
}
