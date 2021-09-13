package gcache

import (
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisOptions interface {
	IsRedisCluster() bool
	RedisOptions() *redis.Options
	RedisClusterOptions() *redis.ClusterOptions
}

type Options struct {
	Addr     string   `json:"addr" yaml:"addr"`
	Addrs    []string `json:"addrs" yaml:"addrs"`
	Username string   `json:"username,omitempty" yaml:"username"`
	Password string   `json:"password" yaml:"password"`
	DB       int      `json:"db,omitempty" yaml:"db"`

	DialTimeout  time.Duration `json:"dial_timeout,omitempty" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout,omitempty" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout,omitempty" yaml:"write_timeout"`
	MaxRetries   int           `json:"max_retries,omitempty" yaml:"max_retries"`
	PoolSize     int           `json:"pool_size,omitempty" yaml:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns,omitempty" yaml:"min_idle_conns"`

	UseTLS bool `json:"use_tls" yaml:"use_tls"`
}

func (c *Options) IsRedisCluster() bool { return len(c.Addrs) > 0 }

// RedisOptions 转换为 redis.Options
func (c *Options) RedisOptions() *redis.Options {
	opt := &redis.Options{
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,

		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		MaxRetries:   c.MaxRetries,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
	}

	if c.UseTLS {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return opt
}

// RedisClusterOptions 转换为 redis.Options
func (c *Options) RedisClusterOptions() *redis.ClusterOptions {
	opt := &redis.ClusterOptions{
		Addrs:    c.Addrs,
		Username: c.Username,
		Password: c.Password,

		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		MaxRetries:   c.MaxRetries,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
	}

	if c.UseTLS {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return opt
}

func DefaultOptions() *Options {
	return &Options{
		Addr: "127.0.0.1:6379",
	}
}

func (o *Options) WithAddr(addr string) *Options {
	o.Addr = addr
	return o
}

func (o *Options) WithUsername(username string) *Options {
	o.Username = username
	return o
}

func (o *Options) WithPassword(password string) *Options {
	o.Password = password
	return o
}
func (o *Options) WithDB(db int) *Options {
	o.DB = db
	return o
}
func (o *Options) WithDialTimeout(dialTimeout time.Duration) *Options {
	o.DialTimeout = dialTimeout
	return o
}
func (o *Options) WithReadTimeout(readTimeout time.Duration) *Options {
	o.ReadTimeout = readTimeout
	return o
}
func (o *Options) WithWriteTimeout(writeTimeout time.Duration) *Options {
	o.WriteTimeout = writeTimeout
	return o
}
func (o *Options) WithMaxRetries(maxRetries int) *Options {
	o.MaxRetries = maxRetries
	return o
}
func (o *Options) WithPoolSize(poolSize int) *Options {
	o.PoolSize = poolSize
	return o
}
func (o *Options) WithMinIdleConns(minIdleConns int) *Options {
	o.MinIdleConns = minIdleConns
	return o
}
func (o *Options) WithUseTLS(useTLS bool) *Options {
	o.UseTLS = useTLS
	return o
}
