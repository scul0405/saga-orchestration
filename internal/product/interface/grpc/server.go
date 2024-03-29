package grpc

import (
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/scul0405/saga-orchestration/cmd/product/config"
	"github.com/scul0405/saga-orchestration/internal/product/app"
	pb "github.com/scul0405/saga-orchestration/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type Server struct {
	Port       string
	productApp        app.ProductApplication
	grpcServer *grpc.Server
	pb.UnimplementedProductServiceServer
}

func NewGRPCServer(config config.GRPC, productApp app.ProductApplication) *Server {
	srv := &Server{
		Port: config.Port,
		productApp:  productApp,
	}

	srv.grpcServer = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: config.MaxConnectionIdle * time.Second,
			MaxConnectionAge:  config.MaxConnectionAge * time.Minute,
			Timeout:           config.Timeout * time.Second,
			Time:              config.Time * time.Second,
		}),
		grpc.ChainUnaryInterceptor(
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	pb.RegisterProductServiceServer(srv.grpcServer, srv)

	reflection.Register(srv.grpcServer)
	return srv
}

func (srv *Server) Run() error {
	addr := "0.0.0.0:" + srv.Port
	log.Println("grpc server listening on ", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if err := srv.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (srv *Server) GracefulStop() {
	srv.grpcServer.GracefulStop()
}
