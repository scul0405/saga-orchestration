package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/scul0405/saga-orchestration/cmd/product/config"
	"github.com/scul0405/saga-orchestration/internal/pkg/cache"
	"github.com/scul0405/saga-orchestration/internal/pkg/grpcconn"
	"github.com/scul0405/saga-orchestration/internal/product/eventhandler"
	"github.com/scul0405/saga-orchestration/internal/product/infrastructure/db/postgres"
	grpcclient "github.com/scul0405/saga-orchestration/internal/product/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/internal/product/interface/grpc"
	"github.com/scul0405/saga-orchestration/internal/product/interface/http"
	"github.com/scul0405/saga-orchestration/internal/product/repository/pgrepo"
	"github.com/scul0405/saga-orchestration/internal/product/repository/proxy"
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

	localCache, err := cache.NewLocalCache(ctx, cfg.LocalCache.ExpirationTime)
	if err != nil {
		apiLogger.Fatal(err)
	}

	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:         cfg.RedisCache.Address,
		Password:      cfg.RedisCache.Password,
		PoolSize:      cfg.RedisCache.PoolSize,
		MaxRetries:    cfg.RedisCache.MaxRetries,
		ReadOnly:      true,
		RouteRandomly: true,
	})

	err = redisClient.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		apiLogger.Fatal(err)
	}

	redisCache := cache.NewRedisCache(redisClient, time.Duration(cfg.RedisCache.ExpirationTime)*time.Second)

	// create repositories
	productPgRepo := pgrepo.NewProductRepository(psqlDB)
	productRepo, err := proxy.NewProductRepository(productPgRepo, localCache, redisCache, apiLogger)
	if err != nil {
		apiLogger.Fatal(err)
	}

	categoryPgRepo := pgrepo.NewCategoryRepositoryImpl(psqlDB)

	// create sony flake
	sf, err := sonyflake.NewSonyFlake()
	if err != nil {
		apiLogger.Fatal(err)
	}

	// create services
	productSvc := service.NewProductService(sf, apiLogger, productRepo)
	categorySvc := service.NewCategoryService(sf, apiLogger, categoryPgRepo)

	authConn, err := grpcconn.NewGRPCClientConn(cfg.RpcEnpoints.AuthSvc)
	if err != nil {
		apiLogger.Fatal(err)
	}
	authSvc := grpcclient.NewAuthService(authConn)

	// create http server
	engine := http.NewEngine(cfg.HTTP)
	router := http.NewRouter(productSvc, categorySvc, authSvc)
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
