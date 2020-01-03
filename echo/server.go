package main

import (
	"context"
	"fmt"
	"github.com/panshul007/grpc-playground/api/echo"
	"google.golang.org/grpc"
	"net"
)

type EchoServer struct{}

func (e *EchoServer) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Response: "My Echo: " + req.Message,
	}, nil
}

func main() {
	fmt.Println("Hello gRPC Playground!!")
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	echoServer := &EchoServer{}
	echo.RegisterEchoServerServer(server, echoServer)

	fmt.Println("Serving echo server at 8080")
	err = server.Serve(lst)
	if err != nil {
		panic(err)
	}
}
