package postgresql

import (
	"ChainServer/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/sea-project/go-logger"
	"time"
)

var (
	Db *gorm.DB
)

func PgOrmOpen() {
	var err error
	host := config.Conf.Pgsql.Host
	user := config.Conf.Pgsql.User
	password := config.Conf.Pgsql.PassWord
	dbname := config.Conf.Pgsql.DbName
	port := config.Conf.Pgsql.Port
	pgInfo := fmt.Sprintf("host=%v port=%v user=%v  dbname=%v password=%v sslmode=disable", host, port, user, dbname, password)
	logger.Info(pgInfo)
	Db, err = gorm.Open("postgres", pgInfo)
	sqlDB := Db.DB()
	sqlDB.SetMaxIdleConns(5)            //空闲连接数
	sqlDB.SetMaxOpenConns(50)           //最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) //过期时间
	Db.SingularTable(true)              // 全局禁用复数表名
	Db.LogMode(true)                    // 显示sql语句调试
	if err != nil {
		panic(err)
	}
}
