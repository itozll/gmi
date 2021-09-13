package gmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	coll *mongo.Collection
}

var (
	FindOne          = options.FindOne
	Find             = options.Find
	InsertOne        = options.InsertOne
	InsertMany       = options.InsertMany
	Update           = options.Update
	FindOneAndUpdate = options.FindOneAndUpdate
	Count            = options.Count
)

func New(tbl, dbName string) (*Model, error) {
	cli, err := Get(dbName)
	if err != nil {
		return nil, err
	}

	return cli.GetModel(tbl), nil
}

func New2(tbl, dbName string) *Model {
	cli, err := Get(dbName)
	if err != nil {
		panic(err)
	}

	return cli.GetModel(tbl)
}

// Clone clone self
func (c *Model) Clone() *Model {
	coll, _ := c.coll.Clone()
	return &Model{coll: coll}
}

func (c *Model) Collection() *mongo.Collection { return c.coll }
func (c *Model) Table() string                 { return c.coll.Name() }
func (c *Model) IsMongo() bool                 { return true }

func (c *Model) GetByObjectId(ctx context.Context, out interface{}, objectId string) (interface{}, error) {
	objId, err := primitive.ObjectIDFromHex(objectId)
	if err != nil {
		return nil, err
	}

	return c.Get(ctx, out, bson.M{"_id": objId})
}

func (c *Model) GetById(ctx context.Context, out interface{}, id interface{}) (interface{}, error) {
	return c.Get(ctx, out, bson.M{"id": id})
}

// Get one by id
func (c *Model) Get(ctx context.Context, out interface{}, filter interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	err := c.coll.FindOne(ctx, filter, opts...).Decode(out)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Gets more by filter
func (c *Model) Gets(ctx context.Context, out interface{}, filter interface{}, opts ...*options.FindOptions) error {
	cur, err := c.coll.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return cur.All(ctx, out)
}

// Insert document
func (c *Model) Insert(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.coll.InsertOne(ctx, document, opts...)
}

// InsertMany document
func (c *Model) InsertMany(ctx context.Context, document []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return c.coll.InsertMany(ctx, document, opts...)
}

// Update one
func (c *Model) Update(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.coll.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany one
func (c *Model) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.coll.UpdateMany(ctx, filter, update, opts...)
}

func (c *Model) InsertOrUpdate(ctx context.Context, out interface{}, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (interface{}, error) {
	opts = append(opts, options.FindOneAndUpdate().SetUpsert(true))
	err := c.coll.FindOneAndUpdate(ctx, filter, update, opts...).Decode(out)
	if err == mongo.ErrNoDocuments {
		return out, nil
	}

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Model) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.coll.CountDocuments(ctx, filter, opts...)
}
