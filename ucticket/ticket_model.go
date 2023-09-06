package ucticket

type ITicketStore interface {
	LoadIDSegment(bizTag string) (ticket *TicketSegment, err error)
	LoadIDSegmentWithNum(bizTag string, num int64) (ticket *TicketSegment, err error)
	InitScope(bizTag string, step int, maxId int64) (err error)
}

type TicketSegment struct {
	BizTag string `bson:"biz_tag" orm:"column(biz_tag);pk"`
	MaxId  int64  `bson:"max_id" orm:"column(max_id)"`
	Step   int    `bson:"step" orm:"column(step)"`
}

func NewTicketStore(ticketUseMode string) ITicketStore {
	if ticketUseMode == TICKET_MODE_MONGO {
		return new(MongoTicketStore)
	} else {
		return new(SQLTicketStore)
	}
}
