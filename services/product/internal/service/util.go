package service

import "context"

func encodeAuthRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeAuthResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	return grpcReply, nil
}
