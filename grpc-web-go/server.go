package main

import (
	"fmt"
	"github.com/oneeyedsunday/go_playground/grpc-web-go/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 14586))

	if err != nil {
		log.Fatalf("failed to listen:  %v\n", err)
	}

	s := new(server.Server)
	grpcServer := grpc.NewServer()

	server.RegisterTodoServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server grpc: %v\n", err)
		return
	}

	log.Printf("server started successfully\n")
}
