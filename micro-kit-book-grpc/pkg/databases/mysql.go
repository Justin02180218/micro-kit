package databases

import (
	"com/justin/micro/kit/pkg/configs"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitMySql(cfg *configs.MySQLConfig) (err error) {
	connUrl := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Db)
	DB, err = gorm.Open("mysql", connUrl)
	if err != nil {
		log.Println(err)
		return
	}

	err = DB.DB().Ping()
	if err != nil {
		log.Println(err)
		return
	}

	if cfg.Debug {
		DB = DB.Debug()
	}
	DB.DB().SetConnMaxLifetime(time.Minute * 10)
	DB.SingularTable(true)
	return
}
