package main

import (
	"fmt"
	"time"
)

func generate(out chan [2]int) {
	for x := 1; x <= 10; x++ {
		for y := 1; y <= 10; y++ {
			t := [2]int{x, y}
			fmt.Printf("generate %dx%d\n ", t[0], t[1])
			out <- t
		}
	}

	fmt.Println("generate end")
}

func multiply(in chan [2]int) {
	for {
		t := <-in
		r := [3]int{t[0], t[1], t[0] * t[1]}
		fmt.Printf("multiply %dx%d=%d\n", r[0], r[1], r[2])
	}
}

func main() {
	in := make(chan [2]int, 10)
	go multiply(in)
	generate(in)
	time.Sleep(10 * time.Second)

}
