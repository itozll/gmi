package mysql_test

import (
	"context"
	"testing"

	"github.com/itozll/gmi/pkg/tool/storage/mysql"
)

type C struct {
	Id       int64 `gorm:"column:id"`
	Clicked  int   `gorm:"column:clicked"`
	Received int   `gorm:"columnd:recevied"`
}

var ctx = context.Background()

func init() {
	mysql.InitByOptions(map[string]*mysql.Options{
		"default": {
			Dsn:   "root:qweASD12#@/push_statistics?charset=utf8&parseTime=True",
			Debug: true,
		},
	})
}

func TestGetByID(t *testing.T) {
	m, err := mysql.New("summary1", "default")
	if err != nil {
		t.Fatal(err)
	}

	var out = []map[string]interface{}{}

	res, err := m.Gets(ctx, &out, map[string]interface{}{
		"id": []int64{1, 2, 3, 4},
	})
	t.Log(err, res)

	var c = C{Id: 1, Clicked: 101}
	err = m.InsertOrUpdate2(
		ctx,
		c,
		"clicked", "received",
	)

	t.Log(err)
}

func TestTransaction(t *testing.T) {
	m, err := mysql.New("summary1", "default")
	if err != nil {
		t.Fatal(err)
	}

	m.Transaction(func(cli *mysql.Model) error {
		cli.Update(ctx, &C{
			Clicked:  103,
			Received: 203,
		}, map[string]interface{}{
			"id": 1,
		})

		cli.Insert(ctx, []C{
			{Id: 2, Clicked: 103},
			{Id: 3, Clicked: 103},
			{Id: 4, Clicked: 103},
			{Id: 5, Clicked: 103},
			{Id: 6, Clicked: 103},
		})

		return nil
		// return errors.New("error")
	})

	var values []C
	m.Gets(ctx, &values, map[string]interface{}{
		"id": []int64{1, 2, 3, 4, 5, 6},
	})
	t.Log(values)

	var value = C{Id: 1}
	m.Gets(ctx, &value, value)

	t.Log(value)

	var value1 C
	m.Gets(ctx, &value1, map[string]interface{}{"id": 2})

	t.Log(value1)
}
