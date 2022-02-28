package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/pgbytes/grpc-playground/api/go/echo"
	"github.com/pgbytes/grpc-playground/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type EchoServer struct{}

func (e *EchoServer) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Response: fmt.Sprintf("My Echo: %s: %+v", req.Message, req.Talk),
	}, nil
}

func valid(auth []string) bool {
	if len(auth) > 1 {
		return false
	}
	token := strings.TrimPrefix(auth[0], "Bearer ")
	return token == "some-super-secret"
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("request intercepted... authenticating...")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	m, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("RPC failed with error: %v \n", err)
	}
	return m, err
}

func main() {
	fmt.Println("Hello gRPC Playground!!")
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// create TLS creds
	cert, err := tls.LoadX509KeyPair(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	if err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)), // enable TLS check.
		grpc.UnaryInterceptor(authInterceptor),              // configure unary methods interceptor to check auth.
	}
	server := grpc.NewServer(opts...)

	echoServer := &EchoServer{}
	echo.RegisterEchoServiceServer(server, echoServer)

	fmt.Println("Serving echo server at 8080")
	err = server.Serve(lst)
	if err != nil {
		panic(err)
	}
}
