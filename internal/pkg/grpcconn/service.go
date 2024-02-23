package grpcconn

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"time"
)

const (
	timeout = 5 * time.Second
)

func NewGRPCClient[GRPCReply any](svcName string, methodName string, conn *GRPCClientConn, reply GRPCReply) endpoint.Endpoint {
	var opts []grpc.ClientOption

	var grpcEndpoint endpoint.Endpoint
	{
		grpcEndpoint = grpc.NewClient(
			conn.Conn,
			svcName,
			methodName,
			encodeGRPCRequest,
			decodeGRPCResponse,
			reply,
			append(opts, grpc.ClientBefore(grpc.SetRequestHeader("Service-Name", svcName)))...,
		).Endpoint()

		grpcEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    svcName,
			Timeout: timeout,
		}))(grpcEndpoint)
	}

	return grpcEndpoint
}
