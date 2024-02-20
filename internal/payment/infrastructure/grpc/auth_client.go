package grpc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/scul0405/saga-orchestration/internal/payment/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/pkg/grpcconn"
	pb "github.com/scul0405/saga-orchestration/proto"
)

type AuthService interface {
	Auth(ctx context.Context, accessToken string) (*valueobject.AuthResponse, error)
}

type authServiceImpl struct {
	auth endpoint.Endpoint
}

func NewAuthService(conn *grpcconn.GRPCClientConn) AuthService {
	authSvc := grpcconn.NewGRPCClient("auth.AuthService", "Auth", conn)

	return &authServiceImpl{
		auth: authSvc,
	}
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
