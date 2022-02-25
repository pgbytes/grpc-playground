package main

import (
	"context"
	"net"
	"testing"

	"github.com/pgbytes/grpc-playground/api/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type EchoTestSuite struct {
	suite.Suite
	listener   net.Listener
	echoServer *grpc.Server
	conn       *grpc.ClientConn
}

func (e *EchoTestSuite) SetupSuite() {
	t := e.T()
	t.Log("creating a server listener at: 8080 ")
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	echoServer := &EchoServer{}
	echo.RegisterEchoServiceServer(server, echoServer)

	t.Log("Serving echo server at 8080")
	go func() {
		err = server.Serve(lst)
		if err != nil {
			panic(err)
		}
	}()
	e.listener = lst
	e.echoServer = server

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		panic(err)
	}
	t.Log("established a connection to localhost:8080")
	e.conn = conn
}

func (e *EchoTestSuite) TearDownSuite() {
	e.T().Log("stopping echo server...")
	e.echoServer.Stop()
	e.T().Log("stopping echo server client connection...")
	err := e.conn.Close()
	if err != nil {
		e.T().Logf("error shutting down echo client: %+v", err)
	}
}

func TestEchoTestSuite(t *testing.T) {
	suite.Run(t, new(EchoTestSuite))
}

func (e *EchoTestSuite) TestEchoServer_Echo() {
	t := e.T()
	client := echo.NewEchoServiceClient(e.conn)
	testCases := []struct {
		description    string
		request        *echo.EchoRequest
		response       *echo.EchoResponse
		expectedStatus error
	}{
		{
			description:    "golden case: no error",
			request:        &echo.EchoRequest{Message: "golden case", Talk: &echo.Message{Text: "golden"}},
			response:       &echo.EchoResponse{},
			expectedStatus: status.New(codes.OK, "").Err(),
		},
		{
			description:    "unauthenticated error",
			request:        &echo.EchoRequest{Message: "unauthenticated request", ErrorType: echo.ErrorType_ERROR_TYPE_UNAUTHENTICATED},
			response:       nil,
			expectedStatus: errUnauthenticated,
		},
		{
			description:    "permission denied error",
			request:        &echo.EchoRequest{Message: "unauthenticated user", ErrorType: echo.ErrorType_ERROR_TYPE_PERMISSION_DENIED},
			response:       nil,
			expectedStatus: errPermissionDenied,
		},
	}
	for _, tc := range testCases {
		t.Logf("testing: %s", tc.description)
		resp, err := client.Echo(context.Background(), tc.request)
		if tc.response != nil {
			assert.NotNil(t, resp)
		} else {
			assert.Nil(t, resp)
		}
		statusErr, isStatus := status.FromError(err)
		require.Truef(t, isStatus, "response error should be from the status package if it is a grpc response")
		expectedStatusErr, _ := status.FromError(tc.expectedStatus)
		assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
	}
}
