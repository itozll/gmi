package ggorm

import (
	"errors"
	"sync"

	"github.com/itozll/gmi/pkg/gconf"
)

var (
	mgr  = map[string]*Client{}
	once sync.Once
	_err error

	ConfigFile = "mysql.yaml"
)

func Init() (err error) {
	once.Do(func() {
		var opts map[string]*Options

		f, err := gconf.AutoLoadYaml(&opts, ConfigFile)
		if err != nil {
			_err = errors.New(f + ": " + err.Error())
			return
		}

		_err = InitByOptions(opts)
	})

	return _err
}

func InitByOptions(opts map[string]*Options) error {
	for key, opt := range opts {
		cli, err := newClient(opt)
		if err != nil {
			return errors.New(key + ": " + err.Error())
		}

		mgr[key] = cli
	}

	return nil
}

func Get(name string) (*Client, error) {
	m, ok := mgr[name]
	if !ok {
		return nil, errors.New("no mysql named: " + name)
	}

	return m, nil
}

func GetModel(table, name string) (*Model, error) {
	return New(table, name)
}
