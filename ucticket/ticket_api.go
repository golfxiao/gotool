package ucticket

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type BiztagType string

func (b BiztagType) String() string {
	return string(b)
}

func (b BiztagType) GetGlobalId() (globalId int64, err error) {
	ticket, ok := ticketObject[b.String()]
	if !ok {
		return 0, fmt.Errorf("Not found biztag: %s", b.String())
	}
	return ticket.Next()
}

func (b BiztagType) GetGlobalIdBatch(num int64) (count []int64, err error) {
	ticket, ok := ticketObject[b.String()]
	if !ok {
		return []int64{}, fmt.Errorf("Not found biztag: %s", b.String())
	}
	return ticket.NextNum(num)
}
