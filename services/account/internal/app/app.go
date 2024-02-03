package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/scul0405/saga-orchestration/services/account/config"
	"github.com/scul0405/saga-orchestration/services/account/internal/infrastructure/db/postgres"
	"github.com/scul0405/saga-orchestration/services/account/internal/infrastructure/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second
)

type app struct {
	log     logger.Logger
	cfg     *config.Config
	pgxConn *pgxpool.Conn
	doneCh  chan struct{}
}

func NewApp(cfg *config.Config, log logger.Logger) *app {
	return &app{
		log:    log,
		cfg:    cfg,
		doneCh: make(chan struct{}),
	}
}

func (a *app) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// connect postgres
	psqlDB, err := postgres.NewPsqlDB(a.cfg)
	if err != nil {
		return err
	}
	defer func() {
		db, err := psqlDB.DB()
		if err = db.Close(); err != nil {
			a.log.Errorf("Close db err: %v", err)
		}
	}()
	// run migration
	a.log.Infof("Run migrations with config: %+v", a.cfg.Migration)
	err = postgres.NewMigrator(psqlDB).Migrate(a.cfg.Migration)
	if err != nil {
		a.log.Errorf("RunMigrations err: %v", err)
		return err
	}
	a.log.Info("Migrations successfully")

	// graceful shutdown
	<-ctx.Done()
	go func() {
		time.Sleep(shutdownTimeout)
		a.log.Infof("Shutdown timeout exceeded, force shutdown")
		a.doneCh <- struct{}{}
	}()

	<-a.doneCh
	a.log.Infof("%s app exited properly", a.cfg.Service.Name)
	return nil
}
