package config

import (
	"errors"
	"github.com/scul0405/saga-orchestration/pkg/appconfig"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	App         appconfig.App
	HTTP        HTTP
	RpcEnpoints RpcEndpoints `mapstructure:"rpcEndpoints"`
}

type Service struct {
	Name         string
	Mode         string
	Debug        bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HTTP struct {
	Port string
	Mode string
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type RpcEndpoints struct {
	AuthSvc    string
	ProductSvc string
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
