package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock sync.Mutex
)

func main() {
	go func() {
		fmt.Println("Pre acquire lock in go routine")
		lock.Lock()
		fmt.Println("Post acquire lock in go routine")
	}()
	fmt.Println("This wil block :) 💣")
	time.Sleep(time.Millisecond * 10)
	fmt.Println("Pre acquire lock in main func")
	lock.Lock()
	fmt.Println("Post acquire lock in main func")
}
