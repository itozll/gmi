package mysql

import (
	"errors"
	"sync"

	"github.com/itozll/ddm/conf"
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

		mgr[key] = cli
	}

	return nil
}

func Get(name string) (*Client, error) {
	m, ok := mgr[name]
	if !ok {
		return nil, errors.New("No mysql named: " + name)
	}

	return m, nil
}

func GetTable(name string, table string) (*Client, error) {
	m, ok := mgr[name]
	if !ok {
		return nil, errors.New("No mysql named: " + name)
	}

	return m.Table(table), nil
}
