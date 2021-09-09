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

	type C struct {
		Id       int64 `gorm:"column:id"`
		Clicked  int   `gorm:"column:clicked"`
		Received int   `gorm:"columnd:recevied"`
	}
	var c = C{Id: 1, Clicked: 101}
	err = m.Table("summary").InsertOrUpdate2(
		ctx,
		c,
		"clicked", "received", "content_type",
	)

	t.Log(err)
}
