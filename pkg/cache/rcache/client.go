package rcache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	isCluster bool

	redis.Cmdable
}

func NewRedis(opt RedisOptions) *Client {
	var cli redis.Cmdable
	isCluster := opt.IsRedisCluster()

	if isCluster {
		cli = redis.NewClusterClient(opt.RedisClusterOptions())
	} else {
		cli = redis.NewClient(opt.RedisOptions())
	}

	return &Client{
		isCluster: isCluster,
		Cmdable:   cli,
	}
}

func NewRedisPing(ctx context.Context, opt RedisOptions) (*Client, error) {
	cli := NewRedis(opt)

	err := cli.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return cli, err
}

func (cli *Client) IsRedisCluster() bool { return cli.isCluster }

func (cli *Client) RedisClient() *redis.Client {
	if cli.isCluster {
		return nil
	}

	return cli.Cmdable.(*redis.Client)
}

func (cli *Client) RedisClusterClient() *redis.ClusterClient {
	if cli.isCluster {
		return cli.Cmdable.(*redis.ClusterClient)
	}

	return nil
}
