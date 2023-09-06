package ucticket

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type SQLTicketStore struct{}

func (this *SQLTicketStore) LoadIDSegment(bizTag string) (segment *TicketSegment, err error) {
	return this.loadIDs(bizTag, cfg.Step)
}

func (this *SQLTicketStore) LoadIDSegmentWithNum(bizTag string, num int64) (segment *TicketSegment, err error) {
	return this.loadIDs(bizTag, int(num))
}

func (this *SQLTicketStore) InitScope(bizTag string, step int, maxId int64) (err error) {
	o := getOrm()

	sql := fmt.Sprintf("INSERT INTO %s(`biz_tag`, `max_id`, `step`) VALUES('%s', '%d', '%d')",
		cfg.TableName, bizTag, maxId, step)

	_, err = o.Raw(sql).Exec()
	if err != nil {
		err = fmt.Errorf("insert sql execute error: %s, %s", err.Error(), sql)
		return
	}
	return
}

func (this *SQLTicketStore) loadIDs(bizTag string, step int) (segment *TicketSegment, err error) {
	success := false
	segment = &TicketSegment{}

	o := getOrm()
	defer func() {
		if success {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()

	o.Begin()
	updateSQL := fmt.Sprintf("UPDATE %s SET max_id=max_id+%d WHERE biz_tag='%s'", cfg.TableName, step, bizTag)
	if _, err = o.Raw(updateSQL).Exec(); err != nil {
		err = fmt.Errorf("update biztag %s segement error:%s", bizTag, err.Error())
		return
	}

	querySQL := fmt.Sprintf("SELECT biz_tag, max_id, step FROM %s WHERE biz_tag='%s'", cfg.TableName, bizTag)
	err = o.Raw(querySQL).QueryRow(segment)
	if err != nil {
		err = fmt.Errorf("get biztag %s segement error:%s", bizTag, err.Error())
		return
	}

	segment.Step = step
	success = true
	return segment, nil
}

func getOrm() orm.Ormer {
	o := orm.NewOrm()
	o.Using(TICKET_DB_ALIASE_NAME)
	return o
}
