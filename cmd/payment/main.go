package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/scul0405/saga-orchestration/cmd/payment/config"
	"github.com/scul0405/saga-orchestration/internal/payment/eventhandler"
	"github.com/scul0405/saga-orchestration/internal/payment/infrastructure/db/postgres"
	"github.com/scul0405/saga-orchestration/internal/payment/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/internal/payment/interface/http"
	"github.com/scul0405/saga-orchestration/internal/payment/repository/pgrepo"
	"github.com/scul0405/saga-orchestration/internal/payment/repository/proxy"
	"github.com/scul0405/saga-orchestration/internal/payment/service"
	"github.com/scul0405/saga-orchestration/internal/pkg/cache"
	"github.com/scul0405/saga-orchestration/internal/pkg/grpcconn"
	kafkaClient "github.com/scul0405/saga-orchestration/pkg/kafka"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/pgconn"
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
	log.Println("Start payment service...")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

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

	// create cache
	localCache, err := cache.NewLocalCache(ctx, cfg.LocalCache.ExpirationTime)
	if err != nil {
		apiLogger.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:       cfg.RedisCache.Address,
		Password:   cfg.RedisCache.Password,
		DB:         cfg.RedisCache.DB,
		PoolSize:   cfg.RedisCache.PoolSize,
		MaxRetries: cfg.RedisCache.MaxRetries,
	})
	redisCache := cache.NewRedisCache(redisClient, time.Duration(cfg.RedisCache.ExpirationTime)*time.Second)

	// create repositories
	paymentPgRepo := pgrepo.NewOrderRepository(psqlDB)
	paymentRepo, err := proxy.NewPaymentRepository(paymentPgRepo, localCache, redisCache, apiLogger)
	if err != nil {
		apiLogger.Fatal(err)
	}

	// connect to other services
	authClientConn, err := grpcconn.NewGRPCClientConn(cfg.RpcEnpoints.AuthSvc)
	if err != nil {
		apiLogger.Fatal(err)
	}
	authSvc := grpc.NewAuthService(authClientConn)

	// create services
	paymentSvc := service.NewPaymentService(apiLogger, paymentRepo)

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

	// create kafka
	producer := kafkaClient.NewProducer(apiLogger, cfg.Kafka.Brokers)
	consumer := kafkaClient.NewConsumerGroup(cfg.Kafka.Brokers, apiLogger)

	// create event handler
	paymentEvHandler := eventhandler.NewEventHandler(cfg, apiLogger, consumer, producer, paymentSvc)

	doneCh := make(chan struct{}) // for graceful shutdown

	// run event handler
	paymentEvHandler.Run(ctx)

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
