package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Listener struct {
	Port 	  string	`fig:"port"`
}

type Listenerer interface {
	Listener() *Listener
}

func NewListenerer(getter kv.Getter) Listenerer {
	return &listenerer{
		getter: getter,
	}
}

type listenerer struct {
	getter kv.Getter
	once   comfig.Once
}

func (l *listenerer) Listener() *Listener {
	return l.once.Do(func() interface{} {
		var config Listener
		err := figure.
			Out(&config).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(l.getter, "listener")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out listener"))
		}


		return &config
	}).(*Listener)
}

