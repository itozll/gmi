package gmongo_test

import (
	"context"
	"testing"

	"github.com/itozll/gmi/pkg/tool/storage/gmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type C struct {
	Id       int64 `gorm:"column:id"`
	Clicked  int   `gorm:"column:clicked"`
	Received int   `gorm:"columnd:recevied"`
}

var ctx = context.Background()

func init() {
	gmongo.InitByOptions(map[string]*gmongo.Options{
		"default": {
			URI:      "mongodb://127.0.0.1:27017",
			Database: "partiko",
		},
	})
}

func TestGetByID(t *testing.T) {
	m, err := gmongo.GetModel("Accounts", "default")
	if err != nil {
		t.Fatal(err)
	}

	var out = map[string]interface{}{}
	_, err = m.GetById(ctx, &out, 105)
	t.Log("GetById", err, out)

	var out1 = map[string]interface{}{}
	_, err = m.GetByObjectId(ctx, &out1, "613f7824590d5a582d1c26ca")
	t.Log("GetByObjectId", err, out1)

	var out0 map[string]interface{}
	res, err := m.InsertOrUpdate(
		ctx,
		&out0,
		bson.M{"id": 97, "name": "test"},
		bson.M{"$set": bson.M{"value": 101}},
		gmongo.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	t.Log("InsertOrUpdate", err, res)
}
