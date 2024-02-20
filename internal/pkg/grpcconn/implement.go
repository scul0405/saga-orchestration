package grpcconn

import (
	"context"
)

func encodeGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeGRPCResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	return grpcReply, nil
}
