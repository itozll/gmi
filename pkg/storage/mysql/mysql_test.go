package mysql

import (
	"context"
	"testing"
)

func TestGetByID(t *testing.T) {
	m, err := New(DefaultOption().
		WithDsn("root:qweASD12#@/push_statistics?charset=utf8&parseTime=True").
		WithDebug(true))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()

	var out = []map[string]interface{}{}

	res, err := m.Table("summary").Gets(ctx, &out, map[string]interface{}{
		"id": []int64{1, 2, 3, 4},
	})
	t.Log(err, res)

	m.Table("summary").InsertOrUpdate(ctx, map[string]interface{}{
		"id": 1,
	}, map[string]interface{}{
		"clicked": 100000,
	})
	t.Log(err, res)

	cnt, err := m.Table("summary").Count(ctx, map[string]interface{}{
		"id": 1,
	})

	t.Log(err, cnt)
}
