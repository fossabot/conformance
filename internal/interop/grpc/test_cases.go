// This contains the test cases from grpc-go interop test_utils.go file,
// https://github.com/grpc/grpc-go/blob/master/interop/test_utils.go
// The test cases have been refactored to be compatible with the standard
// library *testing.T and other implementations of the testing interface.

/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package interopgrpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bufbuild/connect-crosstest/internal/testing"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/benchmark/stats"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	testpb "github.com/bufbuild/connect-crosstest/internal/gen/proto/go/grpc/testing"
	interopconnect "github.com/bufbuild/connect-crosstest/internal/interop/connect"
)

const (
	largeReqSize        = 271828
	largeRespSize       = 314159
	initialMetadataKey  = "x-grpc-test-echo-initial"
	trailingMetadataKey = "x-grpc-test-echo-trailing-bin"
)

var (
	reqSizes  = []int{27182, 8, 1828, 45904}
	respSizes = []int{31415, 9, 2653, 58979}
)

// ClientNewPayload returns a payload of the given type and size.
func ClientNewPayload(t testing.TB, payloadType testpb.PayloadType, size int) (*testpb.Payload, error) {
	t.Helper()
	if size < 0 {
		return nil, fmt.Errorf("Requested a response with invalid length %d", size)
	}
	body := make([]byte, size)
	assert.Equal(t, payloadType, testpb.PayloadType_COMPRESSABLE)
	return &testpb.Payload{
		Type: payloadType,
		Body: body,
	}, nil
}

// DoEmptyUnaryCall performs a unary RPC with empty request and response messages.
func DoEmptyUnaryCall(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	reply, err := client.EmptyCall(context.Background(), &testpb.Empty{}, args...)
	require.NoError(t, err)
	assert.True(t, proto.Equal(&testpb.Empty{}, reply))
	t.Successf("succcessful unary call")
}

// DoLargeUnaryCall performs a unary RPC with large payload in the request and response.
func DoLargeUnaryCall(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	pl, err := ClientNewPayload(t, testpb.PayloadType_COMPRESSABLE, largeReqSize)
	require.NoError(t, err)
	req := &testpb.SimpleRequest{
		ResponseType: testpb.PayloadType_COMPRESSABLE,
		ResponseSize: int32(largeRespSize),
		Payload:      pl,
	}
	reply, err := client.UnaryCall(context.Background(), req, args...)
	require.NoError(t, err)
	assert.Equal(t, reply.GetPayload().GetType(), testpb.PayloadType_COMPRESSABLE)
	assert.Equal(t, len(reply.GetPayload().GetBody()), largeRespSize)
	t.Successf("successful large unary call")
}

// DoClientStreaming performs a client streaming RPC.
func DoClientStreaming(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	stream, err := client.StreamingInputCall(context.Background(), args...)
	require.NoError(t, err)
	var sum int
	for _, s := range reqSizes {
		pl, err := ClientNewPayload(t, testpb.PayloadType_COMPRESSABLE, s)
		require.NoError(t, err)
		req := &testpb.StreamingInputCallRequest{
			Payload: pl,
		}
		require.NoError(t, stream.Send(req))
		sum += s
	}
	reply, err := stream.CloseAndRecv()
	require.NoError(t, err)
	assert.Equal(t, reply.GetAggregatedPayloadSize(), int32(sum))
	t.Successf("successful client streaming test")
}

// DoServerStreaming performs a server streaming RPC.
func DoServerStreaming(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	respParam := make([]*testpb.ResponseParameters, len(respSizes))
	for i, s := range respSizes {
		respParam[i] = &testpb.ResponseParameters{
			Size: int32(s),
		}
	}
	req := &testpb.StreamingOutputCallRequest{
		ResponseType:       testpb.PayloadType_COMPRESSABLE,
		ResponseParameters: respParam,
	}
	stream, err := client.StreamingOutputCall(context.Background(), req, args...)
	require.NoError(t, err)
	var rpcStatus error
	var respCnt int
	var index int
	for {
		reply, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		assert.Equal(t, reply.GetPayload().GetType(), testpb.PayloadType_COMPRESSABLE)
		assert.Equal(t, len(reply.GetPayload().GetBody()), respSizes[index])
		index++
		respCnt++
	}
	assert.Equal(t, rpcStatus, io.EOF)
	assert.Equal(t, respCnt, len(respSizes))
	t.Successf("successful server streaming test")
}

// DoPingPong performs ping-pong style bi-directional streaming RPC.
func DoPingPong(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	stream, err := client.FullDuplexCall(context.Background(), args...)
	require.NoError(t, err)
	var index int
	for index < len(reqSizes) {
		respParam := []*testpb.ResponseParameters{
			{
				Size: int32(respSizes[index]),
			},
		}
		pl, err := ClientNewPayload(t, testpb.PayloadType_COMPRESSABLE, reqSizes[index])
		require.NoError(t, err)
		req := &testpb.StreamingOutputCallRequest{
			ResponseType:       testpb.PayloadType_COMPRESSABLE,
			ResponseParameters: respParam,
			Payload:            pl,
		}
		require.NoError(t, stream.Send(req))
		reply, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, reply.GetPayload().GetType(), testpb.PayloadType_COMPRESSABLE)
		assert.Equal(t, len(reply.GetPayload().GetBody()), respSizes[index])
		index++
	}
	require.NoError(t, stream.CloseSend())
	_, err = stream.Recv()
	assert.Equal(t, err, io.EOF)
	t.Successf("successful ping pong")
}

// DoEmptyStream sets up a bi-directional streaming with zero message.
func DoEmptyStream(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	stream, err := client.FullDuplexCall(context.Background(), args...)
	require.NoError(t, err)
	require.NoError(t, stream.CloseSend())
	_, err = stream.Recv()
	assert.Error(t, err)
	assert.Equal(t, err, io.EOF)
	t.Successf("successful empty stream")
}

// DoTimeoutOnSleepingServer performs an RPC on a sleep server which causes RPC timeout.
func DoTimeoutOnSleepingServer(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	stream, err := client.FullDuplexCall(ctx, args...)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			// This emulates the original test case.
			return
		}
	}
	require.NoError(t, err)
	pl, err := ClientNewPayload(t, testpb.PayloadType_COMPRESSABLE, 27182)
	require.NoError(t, err)
	req := &testpb.StreamingOutputCallRequest{
		ResponseType: testpb.PayloadType_COMPRESSABLE,
		Payload:      pl,
	}
	err = stream.Send(req)
	require.NoError(t, err)
	_, err = stream.Recv()
	assert.Equal(t, status.Code(err), codes.DeadlineExceeded)
	t.Successf("successful timeout on sleep")
}

var testMetadata = metadata.MD{
	"key1": []string{"value1"},
	"key2": []string{"value2"},
}

// DoCancelAfterBegin cancels the RPC after metadata has been sent but before payloads are sent.
func DoCancelAfterBegin(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(context.Background(), testMetadata))
	stream, err := client.StreamingInputCall(ctx, args...)
	require.NoError(t, err)
	cancel()
	_, err = stream.CloseAndRecv()
	assert.Equal(t, status.Code(err), codes.Canceled)
	t.Successf("successful cancel after begin")
}

// DoCancelAfterFirstResponse cancels the RPC after receiving the first message from the server.
func DoCancelAfterFirstResponse(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.FullDuplexCall(ctx, args...)
	require.NoError(t, err)
	respParam := []*testpb.ResponseParameters{
		{
			Size: 31415,
		},
	}
	pl, err := ClientNewPayload(t, testpb.PayloadType_COMPRESSABLE, 27182)
	require.NoError(t, err)
	req := &testpb.StreamingOutputCallRequest{
		ResponseType:       testpb.PayloadType_COMPRESSABLE,
		ResponseParameters: respParam,
		Payload:            pl,
	}
	require.NoError(t, stream.Send(req))
	_, err = stream.Recv()
	require.NoError(t, err)
	cancel()
	_, err = stream.Recv()
	assert.Equal(t, status.Code(err), codes.Canceled)
	t.Successf("successful cancel after first response")
}

var (
	initialMetadataValue  = "test_initial_metadata_value"
	trailingMetadataValue = "\x0a\x0b\x0a\x0b\x0a\x0b"
	customMetadata        = metadata.Pairs(
		initialMetadataKey, initialMetadataValue,
		trailingMetadataKey, trailingMetadataValue,
	)
)

func validateMetadata(t testing.TB, header, trailer metadata.MD) {
	assert.Equal(t, len(header.Get(initialMetadataKey)), 1)
	assert.Equal(t, header[initialMetadataKey][0], initialMetadataValue)
	assert.Equal(t, len(trailer[trailingMetadataKey]), 1)
	assert.Equal(t, trailer[trailingMetadataKey][0], trailingMetadataValue)
}

// DoCustomMetadata checks that metadata is echoed back to the client.
func DoCustomMetadata(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	// Testing with UnaryCall.
	pl, err := ClientNewPayload(t, testpb.PayloadType_COMPRESSABLE, 1)
	require.NoError(t, err)
	req := &testpb.SimpleRequest{
		ResponseType: testpb.PayloadType_COMPRESSABLE,
		ResponseSize: int32(1),
		Payload:      pl,
	}
	ctx := metadata.NewOutgoingContext(context.Background(), customMetadata)
	var header, trailer metadata.MD
	args = append(args, grpc.Header(&header), grpc.Trailer(&trailer))
	reply, err := client.UnaryCall(
		ctx,
		req,
		args...,
	)
	require.NoError(t, err)
	assert.Equal(t, reply.GetPayload().GetType(), testpb.PayloadType_COMPRESSABLE)
	assert.Equal(t, len(reply.GetPayload().GetBody()), 1)
	validateMetadata(t, header, trailer)

	// Testing with FullDuplex.
	stream, err := client.FullDuplexCall(ctx, args...)
	require.NoError(t, err)
	respParam := []*testpb.ResponseParameters{
		{
			Size: 1,
		},
	}
	streamReq := &testpb.StreamingOutputCallRequest{
		ResponseType:       testpb.PayloadType_COMPRESSABLE,
		ResponseParameters: respParam,
		Payload:            pl,
	}
	require.NoError(t, stream.Send(streamReq))
	streamHeader, err := stream.Header()
	require.NoError(t, err)
	_, err = stream.Recv()
	require.NoError(t, err)
	require.NoError(t, stream.CloseSend())
	_, err = stream.Recv()
	assert.Equal(t, err, io.EOF)
	streamTrailer := stream.Trailer()
	validateMetadata(t, streamHeader, streamTrailer)
	t.Successf("successful custom metadata")
}

// DoStatusCodeAndMessage checks that the status code is propagated back to the client.
func DoStatusCodeAndMessage(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	code := int32(codes.Unknown)
	msg := "test status message"
	expectedErr := status.Error(codes.Code(code), msg)
	respStatus := &testpb.EchoStatus{
		Code:    code,
		Message: msg,
	}
	// Test UnaryCall.
	req := &testpb.SimpleRequest{
		ResponseStatus: respStatus,
	}
	_, err := client.UnaryCall(context.Background(), req, args...)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
	// Test FullDuplexCall.
	stream, err := client.FullDuplexCall(context.Background(), args...)
	require.NoError(t, err)
	streamReq := &testpb.StreamingOutputCallRequest{
		ResponseStatus: respStatus,
	}
	require.NoError(t, stream.Send(streamReq))
	require.NoError(t, stream.CloseSend())
	_, err = stream.Recv()
	assert.Equal(t, err.Error(), expectedErr.Error())
	t.Successf("successful status code and message")
}

// DoSpecialStatusMessage verifies Unicode and whitespace is correctly processed
// in status message.
func DoSpecialStatusMessage(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	code := int32(codes.Unknown)
	msg := "\t\ntest with whitespace\r\nand Unicode BMP ☺ and non-BMP 😈\t\n"
	expectedErr := status.Error(codes.Code(code), msg)
	req := &testpb.SimpleRequest{
		ResponseStatus: &testpb.EchoStatus{
			Code:    code,
			Message: msg,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.UnaryCall(ctx, req, args...)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
	t.Successf("successful special status message")
}

// DoUnimplementedService attempts to call a method from an unimplemented service.
func DoUnimplementedService(t testing.TB, client testpb.UnimplementedServiceClient) {
	_, err := client.UnimplementedCall(context.Background(), &testpb.Empty{})
	assert.Equal(t, status.Code(err), codes.Unimplemented)
	t.Successf("successful unimplemented service")
}

// DoUnimplementedMethod attempts to call an unimplemented method.
func DoUnimplementedMethod(t testing.TB, cc *grpc.ClientConn) {
	var req, reply proto.Message
	err := cc.Invoke(context.Background(), "/grpc.testing.TestService/UnimplementedCall", req, reply)
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.Unimplemented)
	t.Successf("successful unimplemented method")
}

func DoFailWithNonASCIIError(t testing.TB, client testpb.TestServiceClient, args ...grpc.CallOption) {
	reply, err := client.FailUnaryCall(
		context.Background(),
		&testpb.SimpleRequest{
			ResponseType: testpb.PayloadType_COMPRESSABLE,
		},
	)
	assert.Nil(t, reply)
	assert.Error(t, err)
	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, s.Code(), codes.ResourceExhausted)
	assert.Equal(t, s.Message(), interopconnect.NonASCIIErrMsg)
	t.Successf("successful fail call with non-ASCII error")
}

func doOneSoakIteration(t testing.TB, ctx context.Context, tc testpb.TestServiceClient, resetChannel bool, serverAddr string, dopts []grpc.DialOption) (latency time.Duration, err error) {
	start := time.Now()
	client := tc
	if resetChannel {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(serverAddr, dopts...)
		if err != nil {
			return time.Nanosecond, err
		}
		defer conn.Close()
		client = testpb.NewTestServiceClient(conn)
	}
	// per test spec, don't include channel shutdown in latency measurement
	defer func() { latency = time.Since(start) }()
	// do a large-unary RPC
	pl, err := ClientNewPayload(t, testpb.PayloadType_COMPRESSABLE, largeReqSize)
	require.NoError(t, err)
	req := &testpb.SimpleRequest{
		ResponseType: testpb.PayloadType_COMPRESSABLE,
		ResponseSize: int32(largeRespSize),
		Payload:      pl,
	}
	var reply *testpb.SimpleResponse
	reply, err = client.UnaryCall(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, reply.GetPayload().GetType(), testpb.PayloadType_COMPRESSABLE)
	assert.Equal(t, len(reply.GetPayload().GetBody()), largeRespSize)
	return
}

// DoSoakTest runs large unary RPCs in a loop for a configurable number of times, with configurable failure thresholds.
// If resetChannel is false, then each RPC will be performed on tc. Otherwise, each RPC will be performed on a new
// stub that is created with the provided server address and dial options.
func DoSoakTest(t testing.TB, client testpb.TestServiceClient, serverAddr string, dopts []grpc.DialOption, resetChannel bool, soakIterations int, maxFailures int, perIterationMaxAcceptableLatency time.Duration, overallDeadline time.Time) {
	ctx, cancel := context.WithDeadline(context.Background(), overallDeadline)
	defer cancel()
	iterationsDone := 0
	totalFailures := 0
	hopts := stats.HistogramOptions{
		NumBuckets:     20,
		GrowthFactor:   1,
		BaseBucketSize: 1,
		MinValue:       0,
	}
	h := stats.NewHistogram(hopts)
	for i := 0; i < soakIterations; i++ {
		if time.Now().After(overallDeadline) {
			break
		}
		iterationsDone++
		latency, err := doOneSoakIteration(t, ctx, client, resetChannel, serverAddr, dopts)
		latencyMs := int64(latency / time.Millisecond)
		h.Add(latencyMs)
		if err != nil {
			totalFailures++
			fmt.Fprintf(os.Stderr, "soak iteration: %d elapsed_ms: %d failed: %s\n", i, latencyMs, err)
			continue
		}
		if latency > perIterationMaxAcceptableLatency {
			totalFailures++
			fmt.Fprintf(os.Stderr, "soak iteration: %d elapsed_ms: %d exceeds max acceptable latency: %d\n", i, latencyMs, perIterationMaxAcceptableLatency.Milliseconds())
			continue
		}
		fmt.Fprintf(os.Stderr, "soak iteration: %d elapsed_ms: %d succeeded\n", i, latencyMs)
	}
	var b bytes.Buffer
	h.Print(&b)
	fmt.Fprintln(os.Stderr, "Histogram of per-iteration latencies in milliseconds:")
	fmt.Fprintln(os.Stderr, b.String())
	fmt.Fprintf(os.Stderr, "soak test ran: %d / %d iterations. total failures: %d. max failures threshold: %d. See breakdown above for which iterations succeeded, failed, and why for more info.\n", iterationsDone, soakIterations, totalFailures, maxFailures)
	assert.True(t, iterationsDone >= soakIterations)
	assert.True(t, totalFailures <= maxFailures)
}