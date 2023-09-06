package ucticket

import (
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type Ticket struct {
	bizTag       string     // 业务名称
	maxId        int64      // 当前这段ID的最大值，小于maxID直接从内存中返回，大于maxID时需要从数据库中加载下一段
	cur          int64      // 已经分配过的ID值，刚加载完一段ID时 cur = maxId - step
	step         int        // ID段的大小，默认50
	lock         sync.Mutex // 锁保护，确保分配ID时线程安全
	next         *Ticket    // 下一段可用ID，用于预加载
	nextLoad     chan bool  // 预加载完成的信号
	nextLoadLock sync.Mutex // 预加载信号判断的保护锁
}

func NewTicket(bizTag string) (*Ticket, error) {
	ticket := new(Ticket)
	ticket.bizTag = bizTag
	segment, err := ticket.getIdSegment(bizTag)
	if err != nil {
		log.Printf("get segment of biztag:%s error:%s", bizTag, err.Error())
		return nil, err
	}

	ticket.next = nil
	ticket.maxId = segment.MaxId
	ticket.step = segment.Step
	ticket.cur = ticket.maxId - int64(ticket.step)

	return ticket, nil
}

func (this *Ticket) Current() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()
	log.Printf("biztag %s current id :%d", this.bizTag, this.cur)
	return this.cur
}

func (this *Ticket) Next() (int64, error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if cfg.UsePreload {
		used := float64(this.cur+int64(this.step)-this.maxId) / float64(this.step)
		// log.Printf("current biztag used :%.3f", used)
		if this.next == nil && this.nextLoad == nil && (used >= cfg.PreloadFactor) {
			// start preload task
			SafeGo(this.preload)
			this.nextLoad = make(chan bool)
		}
	}
	//use next segement from memory or db
	if this.maxId < this.cur+int64(1) {

		if this.nextLoad != nil {
			startTime := time.Now()
			<-this.nextLoad
			log.Printf("wait next segment loading cost:%v", time.Since(startTime).String())
		}

		// when no id can fetch, use or load next segement
		if this.next != nil {
			//load next segement from memory by next segement
			this.maxId = this.next.maxId
			this.cur = this.next.cur
			this.step = this.next.step
			this.next = nil
			log.Printf("load next segement from memory by next segement:%v, biztag:%s ", this, this.bizTag)
		} else {
			segment, err := this.getNextSegment(this.bizTag)
			if err != nil {
				log.Printf("get segment of biztag:%s error:%s", this.bizTag, err.Error())
				return 0, err
			}

			this.maxId = segment.MaxId
			this.step = segment.Step
			this.cur = this.maxId - int64(this.step)
		}

	}
	this.cur++
	// log.Printf("biztag %s next id :%d", this.bizTag, this.cur)
	return this.cur, nil
}

func (this *Ticket) preload() {
	log.Printf("start to load next id segement of biztag:%s ", this.bizTag)
	startTime := time.Now()
	segment, err := this.getNextSegment(this.bizTag)
	if err != nil {
		log.Printf("get segment of biztag:%s error:%s", this.bizTag, err.Error())
		this.nextLoad <- false
		this.nextLoad = nil
		return
	}

	this.next = new(Ticket)
	this.next.bizTag = segment.BizTag
	this.next.maxId = segment.MaxId
	this.next.step = segment.Step
	this.next.cur = segment.MaxId - int64(segment.Step)
	this.next.next = nil
	this.nextLoad <- true
	this.nextLoad = nil
	log.Printf("load next segement use time: %v, results:%v", time.Since(startTime).String(), this.next)
}

func (this *Ticket) getNextSegment(bizTag string) (*TicketSegment, error) {
	return this.getIdSegment(bizTag)
}

func (this *Ticket) getNextSegmentWithNum(bizTag string, num int64) (*TicketSegment, error) {
	ticketStore := NewTicketStore(ticketUseMode)
	segment, err := ticketStore.LoadIDSegmentWithNum(bizTag, num)
	if err != nil {
		return nil, err
	}
	return segment, nil
}

func (this *Ticket) getIdSegment(bizTag string) (*TicketSegment, error) {
	ticketStore := NewTicketStore(ticketUseMode)
	segment, err := ticketStore.LoadIDSegment(bizTag)
	if err != nil {
		// process the case for scope data not init
		log.Printf(err.Error())
		err = ticketStore.InitScope(bizTag, cfg.Step, TICKET_DEFAULT_MAX_ID)
		if err != nil {
			return nil, err
		}
		segment, err = ticketStore.LoadIDSegment(bizTag)
	}
	return segment, err
}

func (this *Ticket) NextNum(num int64) ([]int64, error) {
	ret := []int64{}

	this.lock.Lock()
	defer this.lock.Unlock()

	segment, err := this.getNextSegmentWithNum(this.bizTag, num)
	if err != nil {
		log.Printf("get segment of biztag:%s error:%s", this.bizTag, err.Error())
		return ret, err
	}

	for i := segment.MaxId - num + 1; i <= segment.MaxId; i++ {
		ret = append(ret, i)
	}
	return ret, nil
}

func SafeGo(handler func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("A panic occurred during handler call: %v", p)
				log.Printf("Stack info: \n%s", string(debug.Stack()))
			}
		}()
		handler()
	}()
}
