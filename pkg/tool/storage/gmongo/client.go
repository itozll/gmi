package gmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	opt ClientOptions
	cli *mongo.Client

	database *mongo.Database
}

func newClient(opt ClientOptions) (*Client, error) {
	cli, err := mongo.NewClient(opt.ClientOptions())
	if err != nil {
		return nil, err
	}

	return &Client{
		opt: opt,
		cli: cli,
	}, nil
}

func Connect(ctx context.Context, opt ClientOptions) (*Client, error) {
	cli, err := mongo.Connect(ctx, opt.ClientOptions())
	if err != nil {
		return nil, err
	}

	return &Client{
		opt:      opt,
		cli:      cli,
		database: cli.Database(opt.GetDatabase()),
	}, nil
}

func (m *Client) Client() *mongo.Client { return m.cli }
func (m *Client) Connect(ctx context.Context) error {
	err := m.cli.Connect(ctx)
	if err != nil {
		return err
	}

	m.database = m.cli.Database(m.opt.GetDatabase())
	return nil
}

func (m *Client) Disconnect(ctx context.Context) { m.cli.Disconnect(ctx) }
func (m *Client) GetModel(table string) *Model {
	return &Model{coll: m.database.Collection(table)}
}
