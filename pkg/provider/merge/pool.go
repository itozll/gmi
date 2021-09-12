// 合并功能提供商

package merge

type Pool struct {
	*MergeProvider

	callback func(data interface{})
}

func NewPool(opt *Options, callback func(data interface{})) *Pool {
	if callback == nil {
		panic("fn is null")
	}

	p := &Pool{
		MergeProvider: NewBaseMerge(opt),
		callback:      callback,
	}

	p.Run(p.loop)
	return p
}

func (p *Pool) loop(ch chan interface{}) {
	for data := range ch {
		p.callback(data)
	}
}
