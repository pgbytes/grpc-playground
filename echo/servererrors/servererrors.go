package main

import (
	"context"
	"fmt"
	"net"

	"github.com/pgbytes/grpc-playground/api/go/echo"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnauthenticated  = status.Error(codes.Unauthenticated, "Unauthenticated request")
	errPermissionDenied = status.Error(codes.PermissionDenied, "Unauthenticated User")
	errDefault          = status.Error(codes.Unimplemented, "error type handler not implemented")
)

type EchoServer struct{}

func (e *EchoServer) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	if req.GetErrorType() != echo.ErrorType_ERROR_TYPE_UNSPECIFIED {
		return nil, generateErrorDetail(req)
	}
	return &echo.EchoResponse{
		Response: fmt.Sprintf("My Echo: %s: %+v", req.Message, req.Talk),
	}, nil
}

func generateErrorDetail(req *echo.EchoRequest) error {
	switch req.GetErrorType() {
	case echo.ErrorType_ERROR_TYPE_UNAUTHENTICATED:
		return errUnauthenticated
	case echo.ErrorType_ERROR_TYPE_PERMISSION_DENIED:
		return errPermissionDenied
	case echo.ErrorType_ERROR_TYPE_BAD_REQUEST:
		return badRequestError()
	default:
		return errDefault
	}
}

func badRequestError() error {
	st := status.New(codes.InvalidArgument, "bad request")
	st, err := st.WithDetails(&errdetails.BadRequest_FieldViolation{Field: "email", Description: "invalid format"})
	if err != nil {
		return err
	}
	return st.Err()
}

func main() {
	fmt.Println("welcome to grpc-playground for error details..!!")
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	echoServer := &EchoServer{}
	echo.RegisterEchoServiceServer(server, echoServer)

	fmt.Println("Serving echo server at 8080")
	err = server.Serve(lst)
	if err != nil {
		panic(err)
	}
}
