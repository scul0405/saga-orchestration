package main

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/payment/config"
	"github.com/scul0405/saga-orchestration/services/payment/internal/infrastructure/db/postgres"
	"github.com/scul0405/saga-orchestration/services/payment/internal/infrastructure/grpc/auth"
	"github.com/scul0405/saga-orchestration/services/payment/internal/infrastructure/logger"
	"github.com/scul0405/saga-orchestration/services/payment/internal/interface/http"
	"github.com/scul0405/saga-orchestration/services/payment/internal/repository/pg_repo"
	"github.com/scul0405/saga-orchestration/services/payment/internal/service"
	"github.com/scul0405/saga-orchestration/services/payment/pkg"
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
	log.Println("Start order service...")

	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	apiLogger := logger.NewApiLogger(cfg)
	apiLogger.InitLogger()
	apiLogger.Infof("Service Name: %s, LogLevel: %s, Mode: %s", cfg.Service.Name, cfg.Logger.Level, cfg.Service.Mode)

	// connect postgres
	psqlDB, err := postgres.NewPsqlDB(cfg)
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
	orderRepo := pg_repo.NewOrderRepository(psqlDB)

	// create sony flake
	sf, err := pkg.NewSonyFlake()
	if err != nil {
		apiLogger.Fatal(err)
	}

	authConn, err := auth.NewAuthConn(cfg)
	if err != nil {
		apiLogger.Fatal(err)
	}

	// create services
	paymentSvc := service.NewPaymentService(sf, apiLogger, orderRepo)

	authSvc := auth.NewAuthService(cfg, authConn)

	// create http server
	engine := http.NewEngine(cfg.HTTP)
	router := http.NewRouter(paymentSvc, authSvc)
	httpServer := http.NewHTTPServer(cfg.HTTP, apiLogger, engine, router)

	// run http server
	go func() {
		if err := httpServer.Run(); err != nil {
			apiLogger.Fatalf("Run http server err: %v", err)
		}
	}()

	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

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
}
