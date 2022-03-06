package reflection

import (
	"testing"

	"github.com/pgbytes/grpc-playground/api/go/echo"
	"github.com/stretchr/testify/require"
)

func TestExtractByReflection(t *testing.T) {
	type fakeRequestWithMessage struct {
		Message string
	}
	type fakeRequestWithoutAnyMessage struct {
		RequestID string
	}
	something := "fake-request-non-struct"
	testCases := []struct {
		description          string
		request              interface{}
		expectedMessage      string
		expectedMessageFound bool
	}{
		{
			description:          "golden case: message field found",
			request:              &fakeRequestWithMessage{Message: "3sm5akzqp2u0"},
			expectedMessage:      "3sm5akzqp2u0",
			expectedMessageFound: true,
		},
		{
			description:          "failure case: message field found but empty",
			request:              &fakeRequestWithMessage{Message: ""},
			expectedMessage:      "",
			expectedMessageFound: false,
		},
		{
			description:          "failure case: message field not found",
			request:              &fakeRequestWithoutAnyMessage{RequestID: "foobar-request-id"},
			expectedMessage:      "",
			expectedMessageFound: false,
		},
		{
			description:          "failure case: not a struct",
			request:              &something,
			expectedMessage:      "",
			expectedMessageFound: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualMessage, actualFound := extractByReflection(tc.request)
			require.Equal(t, tc.expectedMessage, actualMessage)
			require.Equal(t, tc.expectedMessageFound, actualFound)
		})
	}
}

func TestSetByReflection(t *testing.T) {
	type fakeRequestWithMessage struct {
		Message string
	}
	type fakeRequestWithIntMessage struct {
		Message int64
	}
	type fakeRequestWithUnexportedMessage struct {
		message int64
	}
	type fakeRequestWithoutAnyMessage struct {
		RequestID string
	}
	something := "fake-request-non-struct"
	testCases := []struct {
		description    string
		request        interface{}
		requestMessage string
		expectedErr    error
	}{
		{
			description:    "golden case: message set correctly",
			request:        &fakeRequestWithMessage{},
			requestMessage: "test-message-to-set",
			expectedErr:    nil,
		},
		{
			description:    "failure case: request not a struct",
			request:        &something,
			requestMessage: "test-message-to-set",
			expectedErr:    errNotAStruct,
		},
		{
			description:    "failure case: field not found in request",
			request:        &fakeRequestWithoutAnyMessage{},
			requestMessage: "test-message-to-set",
			expectedErr:    errFieldNotFound,
		},
		{
			description:    "failure case: request object does not have field of type string",
			request:        &fakeRequestWithIntMessage{},
			requestMessage: "test-message-to-set",
			expectedErr:    errFieldNotOfTypeString,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualErr := setByReflection(tc.request, tc.requestMessage)
			require.Equal(t, tc.expectedErr, actualErr)
			if tc.expectedErr == nil {
				req, ok := tc.request.(*fakeRequestWithMessage)
				require.True(t, ok)
				require.Equal(t, tc.requestMessage, req.Message)
			}
		})
	}
}

func TestExtractByProtoReflection(t *testing.T) {
	something := "fake-request-non-struct"
	testCases := []struct {
		description          string
		request              interface{}
		expectedMessage      string
		expectedMessageFound bool
	}{
		{
			description:          "golden case: message field found",
			request:              &echo.EchoRequest{Message: "3sm5akzqp2u0"},
			expectedMessage:      "3sm5akzqp2u0",
			expectedMessageFound: true,
		},
		{
			description:          "failure case: message field found but empty",
			request:              &echo.EchoRequest{Message: ""},
			expectedMessage:      "",
			expectedMessageFound: false,
		},
		{
			description:          "failure case: message field not found",
			request:              &echo.Message{Text: "foobar-request-id"},
			expectedMessage:      "",
			expectedMessageFound: false,
		},
		{
			description:          "failure case: not a struct",
			request:              &something,
			expectedMessage:      "",
			expectedMessageFound: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualMessage, actualFound := extractByProtoReflection(tc.request)
			require.Equal(t, tc.expectedMessage, actualMessage)
			require.Equal(t, tc.expectedMessageFound, actualFound)
		})
	}
}

func TestSetByProtoReflection(t *testing.T) {
	type fakeRequestWithMessage struct {
		Message string
	}
	testCases := []struct {
		description    string
		request        interface{}
		requestMessage string
		expectedErr    error
	}{
		{
			description:    "golden case: message set correctly",
			request:        &echo.EchoRequest{},
			requestMessage: "test-message-to-set",
			expectedErr:    nil,
		},
		{
			description:    "failure case: request not a proto message",
			request:        &fakeRequestWithMessage{},
			requestMessage: "test-message-to-set",
			expectedErr:    errNotAProto,
		},
		{
			description:    "failure case: field not found in request",
			request:        &echo.Message{},
			requestMessage: "test-message-to-set",
			expectedErr:    errFieldNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualErr := setByProtoReflection(tc.request, tc.requestMessage)
			require.Equal(t, tc.expectedErr, actualErr)
			if tc.expectedErr == nil {
				req, ok := tc.request.(*echo.EchoRequest)
				require.True(t, ok)
				require.Equal(t, tc.requestMessage, req.Message)
			}
		})
	}
}

// compiler optimization for benchmark tests
// according to this: https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
var benchMarkResultProto string
var benchMarkResultProtoBool bool

func BenchmarkExtractByProtoReflection(b *testing.B) {
	request := &echo.EchoRequest{Message: "3sm5akzqp2u0"}
	var actualMessage string
	var actualFound bool
	for i := 0; i < b.N; i++ {
		actualMessage, actualFound = extractByProtoReflection(request)
	}
	benchMarkResultProto = actualMessage
	benchMarkResultProtoBool = actualFound
}

func BenchmarkExtractByReflection(b *testing.B) {
	request := &echo.EchoRequest{Message: "3sm5akzqp2u0"}
	var actualMessage string
	var actualFound bool
	for i := 0; i < b.N; i++ {
		actualMessage, actualFound = extractByReflection(request)
	}
	benchMarkResultProto = actualMessage
	benchMarkResultProtoBool = actualFound
}
