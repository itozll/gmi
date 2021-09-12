package rcache

import (
	"context"
	"time"
)

type (
	locker struct {
		r *Client
	}

	Locker interface {
		TryLock(context.Context, string, time.Duration) bool
		Unlock(context.Context, string)

		// slow
		Lock(ctx context.Context, key string, expire time.Duration, timeout ...time.Duration) bool
	}

	LockerOptions struct {
		Server  string `json:"server,omitempty" yaml:"server" toml:"server"`
		Default bool   `json:"default" yaml:"default" toml:"default"`
	}
)

func NewLocker(r *Client) Locker {
	return &locker{r: r}
}

var _defaultLockTimeout = 1 * time.Second

func (l *locker) Lock(ctx context.Context, key string, expire time.Duration, timeout ...time.Duration) bool {
	_timeout := _defaultLockTimeout
	if len(timeout) > 0 {
		_timeout = timeout[0]
	}

	_ctx, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	ch := make(chan bool)

	go func() {
		for {
			select {
			case <-_ctx.Done():
				ch <- false
				return
			default:
				res, _ := l.r.SetNX(ctx, key, "24", expire).Result()
				if !res {
					time.Sleep(0)
					continue
				}

				ch <- true
				return
			}
		}
	}()

	return <-ch
}

func (l *locker) TryLock(ctx context.Context, key string, expire time.Duration) bool {
	res, _ := l.r.SetNX(ctx, key, "24", expire).Result()
	return res
}

func (l *locker) Unlock(ctx context.Context, key string) {
	retries := 0

retries:
	err := l.r.Del(ctx, key).Err()

	if err != nil {
		retries++
		if retries < 3 {
			goto retries
		}
	}
}
