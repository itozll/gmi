// 并发功能提供商

package merge

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrStoped  = errors.New("stoped")
	ErrTimeout = errors.New("timeout")
)

type MergeProvider struct {
	opt *Options

	wg sync.WaitGroup

	chanStop     chan struct{}
	chanReceived chan interface{}
}

func NewBaseMerge(opt *Options) *MergeProvider {
	opt = opt.init()

	p := &MergeProvider{
		opt: opt,

		chanStop:     make(chan struct{}),
		chanReceived: make(chan interface{}, opt.QueueSize),
	}

	return p
}

func (p *MergeProvider) Stop() {
	close(p.chanStop)
	time.AfterFunc(time.Second, func() {
		close(p.chanReceived)
	})
}
func (p *MergeProvider) Wait() { p.wg.Wait() }

func (p *MergeProvider) Push(data interface{}) error {
	select {
	case <-p.chanStop:
		return ErrStoped

	default:
	}

	to := time.NewTimer(p.opt.Timeout)

	select {
	case <-p.chanStop:
		return ErrStoped

	case <-to.C:
		return ErrTimeout

	case p.chanReceived <- data:
	}

	return nil
}

func (p *MergeProvider) Run(fn func(ch chan interface{})) {
	p.wg.Add(int(p.opt.Processor))
	for i := 0; i < int(p.opt.Processor); i++ {
		go func() {
			defer p.wg.Done()
			fn(p.chanReceived)
		}()
	}
}
