package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	opt ClientOptions
	cli *mongo.Client

	database *mongo.Database
}

func New(opt ClientOptions) (*Mongo, error) {
	cli, err := mongo.NewClient(opt.ClientOptions())
	if err != nil {
		return nil, err
	}

	return &Mongo{
		opt: opt,
		cli: cli,
	}, nil
}

func Connect(ctx context.Context, opt ClientOptions) (*Mongo, error) {
	cli, err := mongo.Connect(ctx, opt.ClientOptions())
	if err != nil {
		return nil, err
	}

	return &Mongo{
		opt:      opt,
		cli:      cli,
		database: cli.Database(opt.GetDatabase()),
	}, nil
}

func (m *Mongo) Client() *mongo.Client { return m.cli }
func (m *Mongo) Connect(ctx context.Context) error {
	err := m.cli.Connect(ctx)
	if err != nil {
		return err
	}

	m.database = m.cli.Database(m.opt.GetDatabase())
	return nil
}

func (m *Mongo) Disconnect(ctx context.Context) { m.cli.Disconnect(ctx) }
func (m *Mongo) Collection(coll string) *mongo.Collection {
	return m.database.Collection(coll)
}
