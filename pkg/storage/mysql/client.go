package mysql

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	TypeSqlite = "sqlites"
	TypeMysql  = "mysql"
)

var DefaultType = TypeMysql

type Column = clause.Column

type Client struct {
	db *gorm.DB
}

func newClient(opt *Options) (*Client, error) {
	var dialector gorm.Dialector
	retry := false
	typ := opt.Type

retry:
	switch opt.Type {
	case TypeSqlite:
		dialector = sqlite.Open(opt.Dsn)

	case TypeMysql:
		dialector = mysql.Open(opt.Dsn)

	default:
		if !retry {
			opt.Type = DefaultType
			goto retry
		}

		log.Fatalln("err type of mysql: " + typ + "|" + DefaultType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: opt.SkipDefaultTransaction,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: !opt.PluralTable,
		},
	})

	if err != nil {
		return nil, err
	}

	if opt.Debug {
		db = db.Debug()
	}

	return &Client{db: db}, nil
}

func (m *Client) DB() *gorm.DB { return m.db }
