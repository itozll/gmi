package merge_test

import (
	"log"
	"testing"
	"time"

	"github.com/itozll/gmi/pkg/provider/merge"
)

var _ merge.Handle = &xType{}

type xType struct {
	value []int
}

func (x *xType) Merge(data interface{}) {
	dat := data.(*xType)
	x.value = append(x.value, dat.value...)
}

func (x *xType) Length() int {
	return len(x.value)
}

func TestMap(t *testing.T) {
	fn := func(key string, data interface{}) {
		log.Println("key", key, "value", data.(*xType).value)
	}

	p := merge.NewMap(&merge.Options{
		Processor:     10,
		TriggerNumber: 4,
		Interval:      10 * time.Millisecond,
		WaitNumber:    2,
	}, fn)

	p.Push(map[string]interface{}{
		"1": &xType{value: []int{1, 4, 3, 5}},
		"3": &xType{value: []int{3}},
		"4": &xType{value: []int{3}},
	})

	p.Push(map[string]interface{}{
		"2": &xType{value: []int{2, 3, 4}},
		"3": &xType{value: []int{2}},
		"5": &xType{value: []int{2}},
	})

	time.Sleep(5 * time.Second)
	p.Stop()
	p.Wait()
}
