package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"

)

type Config interface {
	Listenerer
	comfig.Logger
}

type config struct {
	Listenerer
	comfig.Logger

	getter kv.Getter
}

func NewConfig(getter kv.Getter) Config {
	return &config{
		Listenerer: NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		getter:     getter,
	}
}

