package gmongo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/itozll/gmi/pkg/conf"
	"go.mongodb.org/mongo-driver/mongo"
)

type tMgr struct {
	once sync.Once
	m    *Mongo
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
	for key, opt := range opts {
		cli, err := New(opt)
		if err != nil {
			return errors.New(key + ": " + err.Error())
		}

		mgr[key] = &tMgr{m: cli}
	}

	return nil
}

func GetMongo(name string) (*Mongo, error) {
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

func GetCollection(name, coll string) (*mongo.Collection, error) {
	m, err := GetMongo(name)
	if err != nil {
		return nil, err
	}

	return m.Collection(coll), nil
}
