package rcache

import (
	"context"
	"errors"
	"sync"

	"github.com/itozll/gmi/pkg/conf"
)

type (
	Node struct {
		Server string `json:"server,omitempty" yaml:"server"`
		Key    string `json:"key,omitempty" yaml:"key"`
		Field  string `json:"field,omitempty" yaml:"field"`
	}
)

var (
	_err error
	once sync.Once
	mgr  = map[string]*Client{}

	ConfigFile = "redis.yaml"
)

func Init() error {
	once.Do(func() {
		var opts map[string]*Options

		f, err := conf.AutoLoadYaml(&opts, ConfigFile)
		if err != nil {
			_err = errors.New(f + ": " + err.Error())
			return
		}

		_err = InitByOptions(opts)
	})

	return _err
}

func InitByOptions(opts map[string]*Options) error {
	ctx := context.Background()

	for key, opt := range opts {
		m, err := NewRedisPing(ctx, opt)
		if err != nil {
			return errors.New(key + ": " + err.Error())
		}

		mgr[key] = m
	}

	return nil
}

func Get(name string) (*Client, error) {
	m, ok := mgr[name]
	if !ok {
		return nil, errors.New("No redis named:" + name)
	}

	return m, nil
}
