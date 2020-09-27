package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

const maxBufferSize = 1024

func server(ctx context.Context, address string) (err error) {
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	defer pc.Close()

	doneChannel := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChannel <- err
				return
			}

			fmt.Printf("packet received: bytes=%d from=%s\n", n, addr.String())

			deadline := time.Now().Add(*timeout)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChannel <- err
				return
			}

			n, err = pc.WriteTo(buffer[:n], addr)
			if err != nil {
				doneChannel <- err
				return
			}

			fmt.Printf("packet written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
		// case err = <-doneCHannel:
	}

	return
}
