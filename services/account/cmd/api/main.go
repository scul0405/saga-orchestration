package main

import (
	"github.com/scul0405/saga-orchestration/services/account/config"
	"github.com/scul0405/saga-orchestration/services/account/internal/infrastructure/logger"
	"log"
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

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("Service Name: %s, LogLevel: %s, Mode: %s", cfg.Service.Name, cfg.Logger.Level, cfg.Service.Mode)
}
