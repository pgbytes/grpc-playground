package main

import (
	"context"
	"fmt"
	"github.com/panshul007/grpc-playground/api/echo"
	"github.com/panshul007/grpc-playground/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type tokenAuth struct {
	token string
}

// required method to implement the grpc/credentials.PerRPCCredentials input interface
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t tokenAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	ctx := context.Background()

	// setup transport credentials for TLS
	transportCreds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "echo.test.youtube.com")
	if err != nil {
		panic(err)
	}

	// setup per request credentials
	perRPCCreds := tokenAuth{token: "some-super-secret"}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithPerRPCCredentials(perRPCCreds),
		// grpc.WithBlock(), // block the caller until connection is established
	}

	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("connection established to echo server at localhost:8080")

	ec := echo.NewEchoServerClient(conn)
	resp, err := ec.Echo(ctx, &echo.EchoRequest{
		Message: "Hello World!",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got response from server: %s \n", resp.Response)
}
