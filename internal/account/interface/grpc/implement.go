package grpc

import (
	"context"
	"fmt"
	"github.com/scul0405/saga-orchestration/internal/account/domain/valueobject"
	pb "github.com/scul0405/saga-orchestration/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) Auth(ctx context.Context, req *pb.AuthPayload) (*pb.AuthResponse, error) {
	authPayload := &valueobject.AuthPayload{
		AccessToken: req.AccessToken,
	}
	authResponse, err := srv.authSvc.Auth(ctx, authPayload)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err),
		)
	}
	return &pb.AuthResponse{
		CustomerId: authResponse.CustomerID,
		Expired:    authResponse.Expired,
	}, nil
}
