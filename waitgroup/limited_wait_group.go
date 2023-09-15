package waitgroup

import "sync"

type LimitedWaitGroup struct {
	wg      sync.WaitGroup
	tickets chan struct{} // limit the max number of waitgroup
}

func NewLimitedWaitGroup(max int) *LimitedWaitGroup {
	if max <= 0 {
		max = 1
	}
	return &LimitedWaitGroup{
		tickets: make(chan struct{}, max),
	}
}

func (this *LimitedWaitGroup) Add(n int) {
	for i := 0; i < n; i++ {
		this.tickets <- struct{}{}
	}
	this.wg.Add(n)
}

func (this *LimitedWaitGroup) Done() {
	<-this.tickets
	this.wg.Done()
}

func (this *LimitedWaitGroup) Wait() {
	this.wg.Wait()
}
