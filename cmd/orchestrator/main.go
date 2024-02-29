package main

import (
	"context"
	"github.com/scul0405/saga-orchestration/cmd/orchestrator/config"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/app"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/eventhandler"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/scul0405/saga-orchestration/pkg/logger"
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
	log.Println("Start orchestrator service...")

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

	// Init kafka
	producer := kafkaClient.NewProducer(apiLogger, cfg.Kafka.Brokers)
	consumer := kafkaClient.NewConsumerGroup(cfg.Kafka.Brokers, apiLogger)

	// Init app
	orchestratorSvc := app.NewApp(apiLogger, producer)

	// Init event handler
	orchestratorEvHanlder := eventhandler.NewEventHandler(cfg, apiLogger, consumer, orchestratorSvc)

	// create graceful shutdown
	doneCh := make(chan struct{}) // for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// run event handler
	orchestratorEvHanlder.Run(ctx)

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(shutdownTimeout)
		apiLogger.Infof("Shutdown timeout exceeded, force shutdown")

		doneCh <- struct{}{}
	}()

	<-doneCh
	apiLogger.Infof("%s app exited properly", cfg.App.Service.Name)
}
