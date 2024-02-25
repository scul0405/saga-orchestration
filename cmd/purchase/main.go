package main

import (
	"context"
	"github.com/scul0405/saga-orchestration/cmd/purchase/config"
	"github.com/scul0405/saga-orchestration/internal/common"
	"github.com/scul0405/saga-orchestration/internal/pkg/grpcconn"
	"github.com/scul0405/saga-orchestration/internal/purchase/eventhandler"
	"github.com/scul0405/saga-orchestration/internal/purchase/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/internal/purchase/interface/http"
	"github.com/scul0405/saga-orchestration/internal/purchase/service"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	log.Println("Start purchase service...")

	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	apiLogger := logger.NewApiLogger(&cfg.App)
	apiLogger.InitLogger()
	apiLogger.Infof("Service Name: %s, LogLevel: %s, Mode: %s", cfg.App.Service.Name, cfg.App.Logger.Level, cfg.App.Service.Mode)

	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// create sony flake
	sf, err := sonyflake.NewSonyFlake()
	if err != nil {
		apiLogger.Fatal(err)
	}

	// Create connection
	productClientConn, err := grpcconn.NewGRPCClientConn(cfg.RpcEnpoints.ProductSvc)
	if err != nil {
		apiLogger.Fatal(err)
	}
	productSvc := grpc.NewProductService(productClientConn)

	authClientConn, err := grpcconn.NewGRPCClientConn(cfg.RpcEnpoints.AuthSvc)
	if err != nil {
		apiLogger.Fatal(err)
	}
	authSvc := grpc.NewAuthService(authClientConn)

	producer := kafkaClient.NewProducer(apiLogger, cfg.Kafka.Brokers, common.PurchaseTopic)

	// Create event publisher
	evPub := eventhandler.NewPurchaseEventHandler(producer)

	purchaseSvc := service.NewPurchaseService(sf, apiLogger, productSvc, evPub)

	// create http server
	engine := http.NewEngine(cfg.HTTP)
	router := http.NewRouter(purchaseSvc, authSvc)
	httpServer := http.NewHTTPServer(cfg.HTTP, apiLogger, engine, router)

	// run http server
	go func() {
		if err := httpServer.Run(); err != nil {
			apiLogger.Fatalf("Run http server err: %v", err)
		}
	}()

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(shutdownTimeout)
		apiLogger.Infof("Shutdown timeout exceeded, force shutdown")

		err = httpServer.GracefulStop(ctx)
		if err != nil {
			apiLogger.Errorf("httpServer.GracefulStop err: %v", err)
		}

		doneCh <- struct{}{}
	}()

	<-doneCh
	apiLogger.Infof("%s app exited properly", cfg.App.Service.Name)
}
