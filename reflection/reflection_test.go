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
