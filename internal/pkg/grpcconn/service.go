package grpcconn

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/sony/gobreaker"
	"time"
)

const (
	timeout = 5 * time.Second
)

func NewGRPCClient(svcName string, methodName string, conn *GRPCClientConn) endpoint.Endpoint {
	var opts []grpc.ClientOption

	var grpcEndpoint endpoint.Endpoint
	{
		grpcEndpoint = grpc.NewClient(
			conn.Conn,
			svcName,
			methodName,
			encodeGRPCRequest,
			decodeGRPCResponse,
			&pb.AuthResponse{},
			append(opts, grpc.ClientBefore(grpc.SetRequestHeader("Service-Name", svcName)))...,
		).Endpoint()

		grpcEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    svcName,
			Timeout: timeout,
		}))(grpcEndpoint)
	}

	return grpcEndpoint
}
