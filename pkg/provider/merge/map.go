// 合并功能提供商

package merge

import (
	"time"
)

type Map struct {
	*MergeProvider

	callback func(key string, data interface{})
}

func NewMap(opt *Options, callback func(key string, data interface{})) *Map {
	if callback == nil {
		panic("fn is null")
	}

	p := &Map{
		MergeProvider: NewBaseMerge(opt),
		callback:      callback,
	}

	p.Run(p.loop)
	return p
}

func (p *Map) PushKV(key string, value Handle) error {
	return p.MergeProvider.Push(map[string]interface{}{
		key: value,
	})
}

func (p *Map) loop(ch chan interface{}) {
	opt := p.opt

	var tc *time.Ticker
	if opt.Interval == 0 {
		tc = &time.Ticker{}
	} else {
		tc = time.NewTicker(opt.Interval)
	}

	var list = map[string]tickObj{}

	for {
		select {
		case data, ok := <-ch:
			if !ok {
				goto end
			}

			for key, dat := range data.(map[string]interface{}) {
				dat := dat.(Handle)

				v, ok := list[key]
				if !ok {
					if dat.Length() >= opt.TriggerNumber {
						p.callback(key, dat)
						continue
					}

					list[key] = tickObj{value: dat}
					continue
				}

				v.value.Merge(dat)

				if v.value.Length() >= opt.TriggerNumber {
					p.callback(key, v.value)
					delete(list, key)
					continue
				}

				v.Reset()
			}

		case <-tc.C:
			for key, v := range list {
				if v.TickCount() < opt.TickCount && v.value.Length() < opt.WaitNumber {
					v.Incr()
					continue
				}

				p.callback(key, v.value)
				delete(list, key)
			}
		}
	}

end:
	for key, v := range list {
		p.callback(key, v.value)
	}
}
