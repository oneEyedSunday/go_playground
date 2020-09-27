package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

const maxBufferSize = 1024

func client(ctx context.Context, address string, reader io.Reader) (err error) {
	resolvedAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		// Yup https://tour.golang.org/basics/7
		return // I think the var its returning is explicit becuase the return type has a variable name attached?
	}

	conn, err := net.DialUDP("udp", nil, resolvedAddr)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChannel := make(chan error, 1)

	go func() {
		n, err := io.Copy(conn, reader)
		if err != nil {
			doneChannel <- err
			return
		}

		fmt.Printf("packet written: bytes=%d\n", n)

		buffer := make([]byte, maxBufferSize)

		// timeouts & deadline
		deadline := time.Now().Add(*timeout)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			doneChannel <- err
			return
		}

		nRead, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChannel <- err
			return
		}

		fmt.Printf("packet received: bytes=%d from=%s\n", nRead, addr.String())
		doneChannel <- nil
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
		// case err = <-doneChannel
	}
}
