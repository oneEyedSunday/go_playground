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
	}
}

func main() {
	c := make(chan int)
	// sending c <- 5
	// receiving x := <- c
	// are both blocing

	for i := 0; i < 4; i++ {
		go (&Worker{id: i}).process(c)
	}

	for {
		c <- rand.Int()
		time.Sleep(time.Millisecond * 50)
	}
}
