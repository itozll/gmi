package mysql

import (
	"context"
	"errors"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	TypeSqlite = "sqlites"
	TypeMysql  = "mysql"
)

var DefaultType = TypeMysql

type Client struct {
	db *gorm.DB
}

func New(opt *Options) (*Client, error) {
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

func (m *Client) Table(tbl string) *Client {
	return &Client{db: m.db.Table(tbl)}
}

func (m *Client) DB() *gorm.DB { return m.db }

func (m *Client) GetByID(ctx context.Context, out interface{}, id interface{}) (interface{}, error) {
	db := m.db.Where("id = ?", id).Scan(out)
	if db.Error != nil {
		return nil, db.Error
	}

	if db.RowsAffected <= 0 {
		return nil, nil
	}

	return out, nil
}

func (m *Client) Get(ctx context.Context, out interface{}, filter interface{}) (interface{}, error) {
	return m.Gets(ctx, out, filter)
}

func (m *Client) Gets(ctx context.Context, out interface{}, filter interface{}) (interface{}, error) {
	db := m.db.Where(filter).Scan(out)
	if db.Error != nil {
		return nil, db.Error
	}

	if db.RowsAffected <= 0 {
		return nil, nil
	}

	return out, nil
}

func (m *Client) Insert(ctx context.Context, document interface{}) error {
	db := m.db.Create(document)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected <= 0 {
		return errors.New("not created")
	}

	return nil
}

func (m *Client) Update(ctx context.Context, filter interface{}, update interface{}) error {
	db := m.db.Where(filter).Updates(update)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected <= 0 {
		return errors.New("not updated")
	}

	return nil
}

func (m *Client) UpdateMany(ctx context.Context, filter interface{}, update interface{}) error {
	return m.Update(ctx, filter, update)
}

func (m *Client) InsertOrUpdate(ctx context.Context, filter interface{}, update interface{}) {
	var v = map[string]interface{}{}
	db := m.db.Where(filter)
	db0 := db.Scan(&v)
	var err error

	if db0.RowsAffected <= 0 {
		err = m.Insert(ctx, update)
		if err == nil {
			return
		}
	}

	db.Updates(update)
}

func (m *Client) InsertOrUpdate2(ctx context.Context, filter interface{}, update interface{}) {
	var v = map[string]interface{}{}
	db := m.db.Where(filter)
	db0 := db.Scan(&v)
	var err error

	if db0.RowsAffected <= 0 {
		err = m.Insert(ctx, update)
		if err == nil {
			return
		}
	}

	db.Updates(update)
}

func (m *Client) Count(ctx context.Context, filter interface{}) (int64, error) {
	var cnt int64
	db := m.db.Where(filter).Count(&cnt)
	return cnt, db.Error
}
