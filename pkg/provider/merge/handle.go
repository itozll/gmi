package merge

type Handle interface {
	Merge(interface{})
	Length() int
}

type ticker struct {
	tickCount int
}

func (opt *ticker) TickCount() int { return opt.tickCount }

func (opt *ticker) Reset() { opt.tickCount = 0 }

func (opt *ticker) Incr() { opt.tickCount++ }

type tickObj struct {
	value Handle
	ticker
}
