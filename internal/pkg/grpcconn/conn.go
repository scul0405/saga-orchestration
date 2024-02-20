package grpcconn

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

type GRPCClientConn struct {
	Conn *grpc.ClientConn
}

func NewGRPCClientConn(endpoint string) (*GRPCClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	retryOpts := []grpc_retry.CallOption{
		// generate waits between 900ms to 1100ms
		grpc_retry.WithBackoff(grpc_retry.BackoffLinearWithJitter(1*time.Second, 0.1)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	}

	conn, err := grpc.DialContext(ctx, endpoint, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &GRPCClientConn{
		Conn: conn,
	}, nil
}
