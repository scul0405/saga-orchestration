package auth

import (
	"context"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/scul0405/saga-orchestration/services/order/domain/valueobject"
)

func encodeAuthRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeAuthResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	return grpcReply, nil
}

func (svc *authServiceImpl) Auth(ctx context.Context, accessToken string) (*valueobject.AuthResponse, error) {
	resp, err := svc.auth(ctx, &pb.AuthPayload{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	authResp := resp.(*pb.AuthResponse)

	return &valueobject.AuthResponse{
		CustomerID: authResp.CustomerId,
		Expired:    authResp.Expired,
	}, nil
}
