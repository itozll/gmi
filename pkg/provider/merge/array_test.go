package merge_test

import (
	"log"
	"testing"
	"time"

	"github.com/itozll/gmi/pkg/provider/merge"
)

var _ merge.Handle = &yType{}

type yType struct {
	value []int
}

func (x *yType) Merge(data interface{}) {
	dat := data.(int)
	x.value = append(x.value, dat)
}

func (x *yType) Length() int {
	return len(x.value)
}

func TestArray(t *testing.T) {
	fn := func(data merge.Handle) {
		log.Println("value", data.(*yType).value)
	}

	p := merge.NewArray(&merge.Options{
		Processor:     10,
		TriggerNumber: 4,
		Interval:      1 * time.Second,
		WaitTimes:     2,
		WaitNumber:    3,
	}, func() merge.Handle { return &yType{} }, fn)

	p.MultiPush(1, 2, 3, 4, 5)
	p.MultiPush(8, 9, 10, 22, 34, 123, 344, 24)

	time.Sleep(5 * time.Second)
	p.Stop()
	p.Wait()
}
