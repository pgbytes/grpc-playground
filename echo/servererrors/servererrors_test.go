package main

import (
	"context"
	"net"
	"reflect"
	"testing"

	"github.com/pgbytes/grpc-playground/api/go/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
		description   string
		request       *echo.EchoRequest
		response      *echo.EchoResponse
		expectedError error
	}{
		{
			description:   "golden case: no error",
			request:       &echo.EchoRequest{Message: "golden case", Talk: &echo.Message{Text: "golden"}},
			response:      &echo.EchoResponse{},
			expectedError: status.New(codes.OK, "").Err(),
		},
		{
			description:   "unauthenticated error",
			request:       &echo.EchoRequest{Message: "unauthenticated request", ErrorType: echo.ErrorType_ERROR_TYPE_UNAUTHENTICATED},
			response:      nil,
			expectedError: errUnauthenticated,
		},
		{
			description:   "permission denied error",
			request:       &echo.EchoRequest{Message: "unauthenticated user", ErrorType: echo.ErrorType_ERROR_TYPE_PERMISSION_DENIED},
			response:      nil,
			expectedError: errPermissionDenied,
		},
		{
			description:   "bad request: invalid email",
			request:       &echo.EchoRequest{Message: "unauthenticated user", ErrorType: echo.ErrorType_ERROR_TYPE_BAD_REQUEST},
			response:      nil,
			expectedError: badRequestError(),
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
		expectedStatusErr, _ := status.FromError(tc.expectedError)
		assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
		switch expectedStatusErr.Code() {
		case codes.InvalidArgument:
			assertBadRequest(t, expectedStatusErr, statusErr)
		default:
			// no special checks needed
		}
	}
}

func assertBadRequest(t *testing.T, expected, actual *status.Status) {
	// we assume in base of status code INVALID_ARGUMENT, error details will always have the message type errdetails.BadRequest_FieldViolation
	// https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto#L169

	// ensure error type
	require.Equal(t, expected.Code(), codes.InvalidArgument)
	require.Equal(t, actual.Code(), codes.InvalidArgument)

	// ensure error details type
	t.Log("ensuring error details type by reflection")
	for _, detail := range expected.Details() {
		require.Equal(t, reflect.TypeOf(detail).Elem().String(), "errdetails.BadRequest_FieldViolation")
	}
	for _, detail := range actual.Details() {
		require.Equal(t, reflect.TypeOf(detail).Elem().String(), "errdetails.BadRequest_FieldViolation")
	}

	// ensure type using type switch
	t.Log("ensuring error details type by type switch")
	for _, detail := range expected.Details() {
		switch detail.(type) {
		case *errdetails.BadRequest_FieldViolation:
			// nothing to do test is fine, just for sake of completeness add a require statement
			require.Equal(t, reflect.TypeOf(detail).Elem().String(), "errdetails.BadRequest_FieldViolation")
		default:
			require.Failf(t, "wrong error detail type: %s, expected was: errdetails.BadRequest_FieldViolation", reflect.TypeOf(detail).Elem().String())
		}
	}
	for _, detail := range actual.Details() {
		switch detail.(type) {
		case *errdetails.BadRequest_FieldViolation:
			// nothing to do test is fine, just for sake of completeness add a require statement
			require.Equal(t, reflect.TypeOf(detail).Elem().String(), "errdetails.BadRequest_FieldViolation")
		default:
			require.Failf(t, "wrong error detail type: %s, expected was: errdetails.BadRequest_FieldViolation", reflect.TypeOf(detail).Elem().String())
		}
	}

	// ensure error details type by type casting
	t.Log("ensuring error details type by type casting")
	for _, detail := range expected.Details() {
		fv, ok := detail.(*errdetails.BadRequest_FieldViolation)
		require.True(t, ok)
		require.NotNil(t, &fv)
	}
	for _, detail := range actual.Details() {
		fv, ok := detail.(*errdetails.BadRequest_FieldViolation)
		require.True(t, ok)
		require.NotNil(t, &fv)
	}
}
