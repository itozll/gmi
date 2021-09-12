package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/itozll/gmi/pkg/tool/cache/rcache"
)

func TestRedis(t *testing.T) {
	err := rcache.Init()
	if err != nil {
		t.Fatal(err.Error())
	}

	cli, err := rcache.Get("default")
	if err != nil {
		t.Fatal(err.Error())
	}

	res, err := cli.LRange(context.Background(), "story-impression:32156080", 0, -1).Result()
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(res)

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			defer wg.Done()
			lock := rcache.NewLocker(cli)
			val := lock.Lock(context.TODO(), "111", 1*time.Second, 5*time.Second)

			t.Log(i, val, time.Now())
		}()
	}

	wg.Wait()
}
