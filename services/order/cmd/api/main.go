package api

import (
	"github.com/scul0405/saga-orchestration/services/order/config"
	"log"
)

func main() {
	log.Println("Start order service...")

	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	_, err = config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
}
