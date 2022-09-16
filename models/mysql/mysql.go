package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sea-project/go-logger"
	"time"
)

var DB *sql.DB

func init() {
	var err error
	// 创建连接池
	//DB, err = sql.Open("mysql", config.Conf.Mysql.SourceName)
	DB, err = sql.Open("mysql", "chain:CaCa@0501@tcp(10.0.0.134:3306)/caca?parseTime=true&charset=utf8&loc=Local")
	if err != nil {
		logger.Error(err.Error())
	}

	// 设置最大空闲连接数
	//DB.SetMaxIdleConns(config.Conf.Mysql.MaxIdleConns)
	//DB.SetMaxOpenConns(config.Conf.Mysql.MaxOpenConns)
	DB.SetMaxIdleConns(20)
	DB.SetMaxOpenConns(20)
	DB.SetConnMaxLifetime(time.Minute * 10)

	// Ping
	err = DB.Ping()
	if err != nil {
		logger.Error("MySQL connection failed.", "err", err.Error())
	} else {
		logger.Info("MySQL connection successful!")
	}
}
