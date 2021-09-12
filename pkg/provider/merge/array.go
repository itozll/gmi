// 合并功能提供商

package merge

import (
	"time"
)

type Array struct {
	*MergeProvider

	fnObj    func() Handle
	callback func(data Handle)
}

func NewArray(opt *Options, fnObj func() Handle, callback func(data Handle)) *Array {
	if callback == nil {
		panic("fn is null")
	}

	p := &Array{
		MergeProvider: NewBaseMerge(opt),
		fnObj:         fnObj,
		callback:      callback,
	}

	p.Run(p.loop)
	return p
}

func (p *Array) MultiPush(data ...interface{}) error {
	var err error
	for _, dat := range data {
		err = p.Push(dat)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Array) loop(ch chan interface{}) {
	opt := p.opt

	var tc *time.Ticker
	if opt.Interval == 0 {
		tc = &time.Ticker{}
	} else {
		tc = time.NewTicker(opt.Interval)
	}

	var (
		obj  = p.fnObj()
		tick = tickObj{}
	)

	for {
		select {
		case dat, ok := <-ch:
			if !ok {
				goto end
			}

			obj.Merge(dat)

			if obj.Length() >= opt.TriggerNumber {
				p.callback(obj)
				obj = p.fnObj()
				continue
			}

		case <-tc.C:
			if obj.Length() > 0 {
				if tick.tickCount < opt.WaitTimes && obj.Length() < opt.WaitNumber {
					tick.Incr()
					continue
				}

				p.callback(obj)
				obj = p.fnObj()
			}
		}
	}

end:
	if obj.Length() > 0 {
		p.callback(obj)
	}
}
