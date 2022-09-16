package main

import (
	"ChainServer/api"
	"ChainServer/api/handler"
	"ChainServer/config"
	pg "ChainServer/models/postgresql"
	"github.com/jinzhu/gorm"
)

func init() {
	// 初始化配置信息
	config.Initialize("./conf/log.json")
	pg.PgOrmOpen()
	//task.CronTimer() // 启动自动检测交易hash
	handler.AntInit()
	handler.AliInit()
}

var (
	Db *gorm.DB
)

func main() {
	// api服务启动
	api.RouterStart()
}
