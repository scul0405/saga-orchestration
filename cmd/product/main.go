package main

import (
	"context"
	"github.com/scul0405/saga-orchestration/cmd/product/config"
	"github.com/scul0405/saga-orchestration/internal/pkg/grpcconn"
	"github.com/scul0405/saga-orchestration/internal/product/eventhandler"
	"github.com/scul0405/saga-orchestration/internal/product/infrastructure/db/postgres"
	grpcclient "github.com/scul0405/saga-orchestration/internal/product/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/internal/product/interface/grpc"
	"github.com/scul0405/saga-orchestration/internal/product/interface/http"
	"github.com/scul0405/saga-orchestration/internal/product/repository/pg_repo"
	"github.com/scul0405/saga-orchestration/internal/product/service"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/pgconn"
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
	log.Println("Start product service...")

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

	// connect postgres
	psqlDB, err := pgconn.NewPsqlDB(cfg.Postgres.DnsURL)
	if err != nil {
		apiLogger.Fatal(err)
	}
	defer func() {
		db, err := psqlDB.DB()
		if err = db.Close(); err != nil {
			apiLogger.Errorf("Close db err: %v", err)
		}
	}()

	// run migration
	apiLogger.Infof("Run migrations with config: %+v", cfg.Migration)
	err = postgres.NewMigrator(psqlDB).Migrate(cfg.Migration)
	if err != nil {
		apiLogger.Errorf("RunMigrations err: %v", err)
		apiLogger.Fatal(err)
	}
	apiLogger.Info("Migrations successfully")

	// create repositories
	productRepo := pg_repo.NewProductRepository(psqlDB)

	// create sony flake
	sf, err := sonyflake.NewSonyFlake()
	if err != nil {
		apiLogger.Fatal(err)
	}

	// create services
	productSvc := service.NewProductService(sf, apiLogger, productRepo)

	authConn, err := grpcconn.NewGRPCClientConn(cfg.RpcEnpoints.AuthSvc)
	if err != nil {
		apiLogger.Fatal(err)
	}
	authSvc := grpcclient.NewAuthService(authConn)

	// create http server
	engine := http.NewEngine(cfg.HTTP)
	router := http.NewRouter(productSvc, authSvc)
	httpServer := http.NewHTTPServer(cfg.HTTP, apiLogger, engine, router)

	// run http server
	go func() {
		if err := httpServer.Run(); err != nil {
			apiLogger.Fatalf("Run http server err: %v", err)
		}
	}()

	// create grpc server
	grpcServer := grpc.NewGRPCServer(cfg.GRPC, productSvc)

	// run grpc server
	go func() {
		if err := grpcServer.Run(); err != nil {
			apiLogger.Fatalf("Run grpc server err: %v", err)
		}
	}()

	// create kafka
	producer := kafkaClient.NewProducer(apiLogger, cfg.Kafka.Brokers)
	consumer := kafkaClient.NewConsumerGroup(cfg.Kafka.Brokers, apiLogger)

	// create event handler
	productEvHandler := eventhandler.NewEventHandler(cfg, apiLogger, consumer, producer, productSvc)

	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// run event handler
	productEvHandler.Run(ctx)

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(shutdownTimeout)
		apiLogger.Infof("Shutdown timeout exceeded, force shutdown")

		err = httpServer.GracefulStop(ctx)
		if err != nil {
			apiLogger.Errorf("httpServer.GracefulStop err: %v", err)
		}

		grpcServer.GracefulStop()

		doneCh <- struct{}{}
	}()

	<-doneCh
	apiLogger.Infof("%s app exited properly", cfg.App.Service.Name)
}
