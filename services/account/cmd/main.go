package main

import (
	"github.com/scul0405/saga-orchestration/services/account/config"
	"log"
)

func main() {
	log.Println("Start account service")

	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	_, err = config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	log.Println("Load config successfully")
}
