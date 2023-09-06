package ucticket

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/astaxie/beego/orm"
)

/*
 * 建表语句
DROP TABLE IF EXISTS `ticket`;
CREATE TABLE `ticket` (
  `biz_tag` varchar(128) NOT NULL DEFAULT '',
  `max_id` bigint(20) NOT NULL DEFAULT '1',
  `step` int(11) DEFAULT NULL,
  `desc` varchar(256) DEFAULT NULL,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`biz_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

// ticket mode define
const (
	TICKET_MODE_SQL   = "sql"   // for mode: sql
	TICKET_MODE_MONGO = "mongo" // for mode: mongo

	TICKET_DB_ALIASE_NAME = "ticketdb" // for mode: sql
	TICKET_DEFAULT_MAX_ID = 1          // default max id
)

var (
	ticketObject  = make(map[string]*Ticket, 0)
	ticketUseMode = ""
	cfg           TicketConfig  // global config
	mongoPool     *mongo.Client // for mode : mongo
)

type TicketConfig struct {
	DataSrc       string   // DB connection string
	ConnCount     int      // DB connection count
	TableName     string   // ticket table name
	Step          int      // id segment size
	ScopeList     []string // biz_tag list
	UsePreload    bool     // whether use preload
	PreloadFactor float64  // preload factor, from 0.0 to 1.0
	DatabaseName  string   // for mode: mongo
}

// ticket mode : sql
func InitTicketDB(config TicketConfig) (err error) {
	if err = checkConfig(&config); err != nil {
		return
	}

	// register ticket db
	err = orm.RegisterDataBase(TICKET_DB_ALIASE_NAME,
		"mysql", config.DataSrc, config.ConnCount, config.ConnCount)
	if err != nil {
		return
	}
	orm.RegisterModel(new(TicketSegment))

	cfg = config
	ticketUseMode = TICKET_MODE_SQL

	// init and cache id for scopelist
	err = initTicketObject(cfg.ScopeList)
	return
}

func InitTicketMongo(config TicketConfig) (err error) {
	if err = checkConfig(&config); err != nil {
		return
	}

	ctx := context.Background()
	mongoPool, err = mongo.Connect(ctx, options.Client().
		ApplyURI(config.DataSrc).
		SetMinPoolSize(uint64(config.ConnCount)).
		SetMaxPoolSize(uint64(config.ConnCount)))
	if err != nil {
		return
	}

	err = mongoPool.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	cfg = config
	ticketUseMode = TICKET_MODE_MONGO

	// init and cache id for scopelist
	err = initTicketObject(cfg.ScopeList)
	return
}

func checkConfig(config *TicketConfig) error {
	if config.DataSrc == "" {
		return errors.New("invalid params of ticket datasrc empty")
	}
	if len(config.ScopeList) == 0 {
		return errors.New("invalid params of ticket scope list empty")
	}
	if config.ConnCount <= 0 {
		config.ConnCount = 10
	}
	if config.TableName == "" {
		config.TableName = "ticket"
	}
	if config.Step <= 0 {
		config.Step = 50
	}
	if config.UsePreload && config.PreloadFactor <= 0.0 {
		config.PreloadFactor = 0.75
	}
	return nil
}

func initTicketObject(scope []string) error {
	for _, bizTag := range scope {
		ticket, err := NewTicket(bizTag)
		if err != nil {
			return err
		}
		if ticket == nil {
			err = fmt.Errorf("init ticket object by biztag :%s is empty", bizTag)
			return err
		}
		ticketObject[bizTag] = ticket
	}

	log.Printf("ticket object :%+v", ticketObject)
	return nil
}
