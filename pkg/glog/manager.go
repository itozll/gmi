package glog

import (
	"errors"
	"sync"

	"github.com/itozll/gmi/pkg/gconf"
)

var (
	once sync.Once
	_err error
	mgr  = map[string]Logger{}

	ConfigFile = "logger.yaml"
)

func Init() error {
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
		l := New(opt)
		mgr[key] = l
		if key == "default" {
			dfl = l
		}
	}

	return nil
}

func Get(name string) Logger {
	m, ok := mgr[name]
	if !ok {
		return dfl
	}

	return m
}
