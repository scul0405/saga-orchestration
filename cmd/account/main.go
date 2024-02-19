package main

import (
	"context"
	"github.com/scul0405/saga-orchestration/cmd/account/config"
	"github.com/scul0405/saga-orchestration/internal/account/infrastructure/db/postgres"
	portgrpc "github.com/scul0405/saga-orchestration/internal/account/ports/grpc"
	porthttp "github.com/scul0405/saga-orchestration/internal/account/ports/http"
	"github.com/scul0405/saga-orchestration/internal/account/repository/postgres_repo"
	customersvc "github.com/scul0405/saga-orchestration/internal/account/service/account"
	authsvc "github.com/scul0405/saga-orchestration/internal/account/service/auth"
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
	log.Println("Start account service")

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
	customerRepo := postgres_repo.NewCustomerRepositoryImpl(psqlDB)
	jwtAuthRepo := postgres_repo.NewJwtAuthRepositoryImpl(psqlDB)

	// create sony flake
	sf, err := sonyflake.NewSonyFlake()
	if err != nil {
		apiLogger.Fatal(err)
	}

	// create services
	customerService := customersvc.NewCustomerService(customerRepo, apiLogger)
	authService := authsvc.NewJWTAuthService(cfg.JWTConfig, jwtAuthRepo, apiLogger, sf)

	// create http server
	engine := porthttp.NewEngine(cfg.HTTP)
	router := porthttp.NewRouter(authService, customerService)
	httpServer := porthttp.NewHTTPServer(cfg.HTTP, apiLogger, engine, router)

	// run http server
	go func() {
		if err := httpServer.Run(); err != nil {
			apiLogger.Fatalf("Run http server err: %v", err)
		}
	}()

	// create grpc server
	grpcServer := portgrpc.NewGRPCServer(cfg.GRPC, authService)

	// run grpc server
	go func() {
		if err := grpcServer.Run(); err != nil {
			apiLogger.Fatalf("Run grpc server err: %v", err)
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

		grpcServer.GracefulStop()
		doneCh <- struct{}{}
	}()

	<-doneCh
	apiLogger.Infof("%s app exited properly", cfg.App.Service.Name)
}
