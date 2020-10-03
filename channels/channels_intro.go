package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Processor interface {
	process(chan int)
}

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
	var _ Processor = &Worker{} // because i defined the method on a pointer
	// var _ Processor = (*Worker)(nil) // interface satisfaction assertion if method defined on value
	c := make(chan int, 0)
	// sending c <- 5
	// receiving x := <- c
	// are both blocing

	for i := 0; i < 4; i++ {
		go (&Worker{id: i}).process(c)
	}

	for {
		select {
		case c <- rand.Int():
			// fmt.Printf("<------ \n")
		case t := <-time.After(time.Millisecond * 700):
			fmt.Println("Timed out at: ", t)
			// default:
			//	fmt.Printf("\tDropped\n")
		}
		// c <- rand.Int()
		// fmt.Printf("Channel length (Buffered) %d capacity is: %d\n", len(c), cap(c))
		time.Sleep(time.Millisecond * 50) // pause exec so user can more easily see flow
	}
}
