package upgrade

import (
	"fmt"
	"gotool/upgrade/models"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/sync/singleflight"

	_ "github.com/go-sql-driver/mysql"
)

var (
	sg = &singleflight.Group{}
)

func Init(confPath string) error {
	beego.AppConfigPath = confPath
	if err := beego.ParseConfig(); err != nil {
		return err
	}
	cons, err := beego.AppConfig.Int("mysql_cons")
	if err != nil {
		return err
	}
	datasrc := beego.AppConfig.String("mysql_datasrc")
	if err := InitConnectionPool(cons, datasrc); err != nil {
		return err
	}

	log.Printf("initialize connection pool success")
	return nil
}

func InitConnectionPool(cons int, datasrc string) error {
	if cons <= 0 || datasrc == "" {
		return fmt.Errorf("params of cons and datasrc can not empty")
	}
	err := orm.RegisterDriver("mysql", orm.DR_MySQL)
	if err != nil {
		return err
	}
	err = orm.RegisterDataBase("default", "mysql", datasrc, cons, cons)
	if err != nil {
		return err
	}
	orm.RegisterModel(new(models.TReleaseNotes), new(models.TReleaseFile))
	return nil
}
