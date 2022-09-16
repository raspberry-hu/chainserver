package model

import (
	. "ChainServer/models/postgresql"
	"fmt"
	"github.com/sea-project/go-logger"
	"time"
)

type TWalletInfo struct {
	Id            string
	WalletAddr    string    `json:"wallet_addr"`    // 钱包地址
	UserName      string    `json:"user_name"`      // 用户名
	ImageUrl      string    `json:"image_url"`      // 头像地址
	BackgroundUrl string    `json:"background_url"` // 背景地址
	EmailAddr     string    `json:"email_addr"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}

func TWalletInfoInsert(t_wallet_info TWalletInfo) bool {
	er := Db.Create(&t_wallet_info)
	flag := Db.NewRecord(t_wallet_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return false
	}
	return true
}

func TWalletInfoFind(sql interface{}) (t_wallet_info []TWalletInfo) {
	Db.Where(sql).Find(&t_wallet_info)
	return t_wallet_info
}

func TaskInfoUpdate(t_wallet_info TWalletInfo, update_info TWalletInfo) bool {
	Db.Model(&t_wallet_info).Updates(update_info)
	return true
}

/*
func TaskInfoUpdate(screen_sql, update_sql interface{}) bool {
	var t_wallet_info TWalletInfo
	Db.Model(&t_wallet_info).Updates(update_sql).Where(screen_sql).Find(&t_wallet_info)
	return true
}
*/

func WalletInfoUpdate(wallet_addr, user_name, image_url string) error {
	var sql string
	if user_name != "" && image_url != "" {
		sql = fmt.Sprintf("UPDATE wallet_info SET user_name = '%s' , image_url = '%s' WHERE wallet_addr = '%s'", user_name, image_url, wallet_addr)
	} else if image_url == "" {
		sql = fmt.Sprintf("UPDATE wallet_info SET  image_url = '%s' WHERE wallet_addr = '%s'", image_url, wallet_addr)
	} else if user_name == "" {
		sql = fmt.Sprintf("UPDATE wallet_info SET  user_name = '%s' WHERE wallet_addr = '%s'", user_name, wallet_addr)
	}
	Db.Exec(sql)
	return nil
}
