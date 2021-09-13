package ggorm

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	tbl string
	db  *gorm.DB
}

func New(tbl, dbName string) (*Model, error) {
	cli, err := Get(dbName)
	if err != nil {
		return nil, err
	}

	return &Model{tbl: tbl, db: cli.db}, nil
}

func New2(tbl, dbName string) *Model {
	cli, err := Get(dbName)
	if err != nil {
		panic(err)
	}

	return &Model{tbl: tbl, db: cli.db}
}

// InsteadOf 替换 *gorm.DB
/**
 * m0 := New("a", "default")
 * m1 := New("b", "default")
 * ...
 * m0.Transaction(func(cli *Model) error {
 * 	cli is "a"
 * 	mm1 := m1.InsteadOf(cli.DB())
 *  now mm1 is "b"
 * })
**/
func (m *Model) InsteadOf(tx *gorm.DB) *Model {
	return &Model{tbl: m.tbl, db: tx}
}

func (m *Model) Table() string { return m.tbl }

func (m *Model) DB() *gorm.DB {
	return m.db.Table(m.tbl)
}

func (m *Model) GetByID(_ context.Context, out interface{}, id interface{}) (interface{}, error) {
	db := m.DB().Where("id = ?", id).Scan(out)
	if db.Error != nil {
		return nil, db.Error
	}

	if db.RowsAffected <= 0 {
		return nil, nil
	}

	return out, nil
}

func (m *Model) Gets(_ context.Context, out interface{}, filter interface{}, args ...interface{}) (interface{}, error) {
	db := m.DB().Where(filter, args...).Scan(out)
	if db.Error != nil {
		return nil, db.Error
	}

	if db.RowsAffected <= 0 {
		return nil, nil
	}

	return out, nil
}

func (m *Model) Insert(_ context.Context, document interface{}) (rowsAffected int64, err error) {
	db := m.DB().Create(document)
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (m *Model) Update(_ context.Context, update interface{}, filter interface{}, args ...interface{}) (rowsAffected int64, err error) {
	db := m.DB().Where(filter, args...).Updates(update)
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

var _columns = []clause.Column{{Name: "id"}}

func (m *Model) InsertOrUpdate(_ context.Context, insert interface{}, update map[string]interface{}) error {
	onConflict := clause.OnConflict{
		Columns: _columns,
	}
	if update == nil {
		onConflict.DoNothing = true
	} else {
		onConflict.DoUpdates = clause.Assignments(update)
	}

	return m.DB().Clauses(onConflict).Create(insert).Error
}

func (m *Model) InsertOrUpdate2(_ context.Context, insert interface{}, updateColumns ...string) error {
	onConflict := clause.OnConflict{
		Columns: _columns,
	}

	if updateColumns == nil {
		onConflict.DoNothing = true
	} else {
		onConflict.DoUpdates = clause.AssignmentColumns(updateColumns)
	}

	return m.DB().Clauses(onConflict).Create(insert).Error
}

func (m *Model) Count(_ context.Context, filter interface{}, args ...interface{}) (int64, error) {
	var cnt int64
	db := m.DB().Where(filter, args...).Count(&cnt)
	return cnt, db.Error
}

func (m *Model) Transaction(fn func(cli *Model) error) error {
	return m.DB().Transaction(func(tx *gorm.DB) error {
		return fn(&Model{tbl: m.tbl, db: tx})
	})
}
