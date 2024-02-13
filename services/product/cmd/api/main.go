package main

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/product/config"
	"github.com/scul0405/saga-orchestration/services/product/internal/infrastructure/db/postgres"
	"github.com/scul0405/saga-orchestration/services/product/internal/infrastructure/logger"
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

	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(shutdownTimeout)
		doneCh <- struct{}{}
	}()

	<-doneCh
}
