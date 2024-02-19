package grpc

import (
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/scul0405/saga-orchestration/cmd/product/config"
	pb "github.com/scul0405/saga-orchestration/proto"
	"github.com/scul0405/saga-orchestration/services/product/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type Server struct {
	Port       string
	app        app.Application
	grpcServer *grpc.Server
	pb.UnimplementedProductServiceServer
}

func NewGRPCServer(config config.GRPC, app app.Application) *Server {
	srv := &Server{
		Port: config.Port,
		app:  app,
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
