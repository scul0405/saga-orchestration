package appconfig

import (
	"time"
)

type App struct {
	Service Service `mapstructure:"service"`
	Logger  Logger  `mapstructure:"logger"`
}

type Service struct {
	Name         string
	Mode         string
	Debug        bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}
