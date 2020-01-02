package main

import (
	"context"
	"fmt"
	"github.com/panshul007/grpc-playground/proto/echo"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ec := echo.NewEchoServerClient(conn)
	resp, err := ec.Echo(ctx, &echo.EchoRequest{
		Message: "Hello World!",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got response from server: %s", resp.Response)
}
