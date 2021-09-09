package mysql

import (
	"context"
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

func (m *Client) GetByID(_ context.Context, out interface{}, id interface{}) (interface{}, error) {
	db := m.db.Where("id = ?", id).Scan(out)
	if db.Error != nil {
		return nil, db.Error
	}

	if db.RowsAffected <= 0 {
		return nil, nil
	}

	return out, nil
}

func (m *Client) Gets(_ context.Context, out interface{}, filter interface{}, args ...interface{}) (interface{}, error) {
	db := m.db.Where(filter, args...).Scan(out)
	if db.Error != nil {
		return nil, db.Error
	}

	if db.RowsAffected <= 0 {
		return nil, nil
	}

	return out, nil
}

func (m *Client) Insert(_ context.Context, document interface{}) (rowsAffected int64, err error) {
	db := m.db.Create(document)
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (m *Client) Update(_ context.Context, update interface{}, filter interface{}, args ...interface{}) (rowsAffected int64, err error) {
	db := m.db.Where(filter, args...).Updates(update)
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

var _columns = []clause.Column{{Name: "id"}}

func (m *Client) InsertOrUpdate(_ context.Context, insert interface{}, update map[string]interface{}) error {
	onConflict := clause.OnConflict{
		Columns: _columns,
	}
	if update == nil {
		onConflict.DoNothing = true
	} else {
		onConflict.DoUpdates = clause.Assignments(update)
	}

	return m.db.Clauses(onConflict).Create(insert).Error
}

func (m *Client) InsertOrUpdate2(_ context.Context, insert interface{}, updateColumns ...string) error {
	onConflict := clause.OnConflict{
		Columns: _columns,
	}

	if updateColumns == nil {
		onConflict.DoNothing = true
	} else {
		onConflict.DoUpdates = clause.AssignmentColumns(updateColumns)
	}

	return m.db.Clauses(onConflict).Create(insert).Error
}

func (m *Client) Count(_ context.Context, filter interface{}, args ...interface{}) (int64, error) {
	var cnt int64
	db := m.db.Where(filter, args...).Count(&cnt)
	return cnt, db.Error
}
