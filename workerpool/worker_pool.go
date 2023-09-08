package workerpool

import (
	"context"
	"errors"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

var (
	ErrParamInvalid   = errors.New("Invalid param of func nil")
	ErrAlreadyRunning = errors.New("A task with the same key is already running")
)

type Ticket int     // 线程资源凭证，申请到资源才能运行任务
type TaskKey string // 任务标识，业务方指定，用于避免相同key重复运行任务
type Task func()    // 待运行任务

type WorkerPool struct {
	sync.RWMutex                      // 对running读写加保护
	tickets      chan Ticket          // 资源池，用于并发任务数量管理，只有申请到ticket才能启动任务，任务运行完会释放ticket
	running      map[TaskKey]struct{} // 任务运行时标记，可选，如果在Run时指定了key, 则一个key同一时间只能运行一个任务
}

// / 构造函数，max参数用于指定能同时运行的最大任务数量
func NewWorkerPool(max int) *WorkerPool {
	p := &WorkerPool{
		tickets: make(chan Ticket, max),
		running: make(map[TaskKey]struct{}, max),
	}
	for i := 0; i < max; i++ {
		p.tickets <- Ticket(i)
	}
	return p
}

func (p *WorkerPool) Run(ctx context.Context, run Task) error {
	return p.RunWithKey(ctx, "", run)
}

func (p *WorkerPool) RunWithKey(ctx context.Context, key TaskKey, run Task) error {
	if run == nil {
		return ErrParamInvalid
	}
	if len(key) > 0 && p.isRunning(key) {
		return ErrAlreadyRunning
	}

	ticket := p.applyTicket()

	// 申请到ticket后需要标记下这个key已经在运行，后续同一个key的任务不再重复申请ticket
	// 这里有一个细节：在打标记时如果标记已经存在，则认为这个key已经有另一个线程申请到ticket并运行，此时需要把当前线程申请到的ticket释放掉；
	if len(key) > 0 && !p.markRunning(key) {
		p.releaseTicket(ticket)
		return ErrAlreadyRunning
	}

	release := func() {
		if len(key) > 0 {
			p.deleteRunning(key)
		}
		p.releaseTicket(ticket)
	}
	e := &worker{run}
	e.Go(release)

	return nil
}

func (p *WorkerPool) applyTicket() Ticket {
	return <-p.tickets
}

func (p *WorkerPool) releaseTicket(t Ticket) {
	p.tickets <- t
}

func (p *WorkerPool) markRunning(key TaskKey) bool {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.running[key]; ok {
		return false
	} else {
		p.running[key] = struct{}{}
		return true
	}
}

func (p *WorkerPool) deleteRunning(key TaskKey) {
	p.Lock()
	delete(p.running, key)
	p.Unlock()
}

func (p *WorkerPool) isRunning(key TaskKey) bool {
	p.RLock()
	_, ok := p.running[key]
	p.RUnlock()
	return ok
}

type worker struct {
	run Task // 待运行任务
}

func (e *worker) Go(release func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("A panic occurred: %v, StackInfo: %s", p, string(debug.Stack()))
			}
			if release != nil {
				release()
			}
		}()
		start := time.Now()
		e.run()
		log.Printf("The time of task running: %s", time.Since(start).String())
	}()
}
