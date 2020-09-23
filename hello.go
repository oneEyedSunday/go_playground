package main

import "fmt"
import "rsc.io/quote"

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(quote.Go())

	var i int = 1
	// switch and fallthrough
	switch i {
	case 0:
		fmt.Println("None")
	case 1:
		fmt.Print("Single")
		fallthrough
	default:
		fmt.Print(" Life...")
	}
}
