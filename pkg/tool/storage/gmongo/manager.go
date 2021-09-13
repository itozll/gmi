package gmongo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/itozll/gmi/pkg/gconf"
)

type tMgr struct {
	once sync.Once
	m    *Client
}

var (
	mgr  = map[string]*tMgr{}
	once sync.Once
	_err error

	ConfigFile = "mongo.yaml"
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
		cli, err := newClient(opt)
		if err != nil {
			return errors.New(key + ": " + err.Error())
		}

		mgr[key] = &tMgr{m: cli}
	}

	return nil
}

func Get(name string) (*Client, error) {
	m, ok := mgr[name]
	if !ok {
		return nil, errors.New("No mongo named: " + name)
	}

	var err error

	m.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()
		err = m.m.Connect(ctx)
	})

	return m.m, err
}

func GetModel(table, name string) (*Model, error) {
	return New(table, name)
}
