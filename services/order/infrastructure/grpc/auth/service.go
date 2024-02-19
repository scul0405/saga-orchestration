package auth

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/scul0405/saga-orchestration/cmd/order/config"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/scul0405/saga-orchestration/services/order/domain/valueobject"
	"github.com/sony/gobreaker"
	"time"
)

const (
	AuthSvcName = "auth.AuthService"
	timeout     = 5 * time.Second
)

type AuthService interface {
	Auth(ctx context.Context, accessToken string) (*valueobject.AuthResponse, error)
}

type authServiceImpl struct {
	auth endpoint.Endpoint
}

func NewAuthService(cfg *config.Config, conn *AuthConn) AuthService {

	var opts []grpc.ClientOption

	var authEndpoint endpoint.Endpoint

	{

		authEndpoint = grpc.NewClient(
			conn.Conn,
			AuthSvcName,
			"Auth",
			encodeAuthRequest,
			decodeAuthResponse,
			&pb.AuthResponse{},
			append(opts, grpc.ClientBefore(grpc.SetRequestHeader("Service-Name", AuthSvcName)))...,
		).Endpoint()

		authEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    AuthSvcName,
			Timeout: timeout,
		}))(authEndpoint)
	}

	return &authServiceImpl{
		auth: authEndpoint,
	}
}
