package models

import (
	. "ChainServer/models/postgresql"
	"fmt"
	"time"

	"github.com/sea-project/go-logger"
)

type TWalletInfo struct {
	Id         int
	WalletAddr string
	UserName   string
	UserDesc   string
	ImageUrl   string
	BannerUrl  string
	EmailAddr  string
	CreateTime time.Time
	UpdateTime time.Time
}

func WalletInfoInsert(wallet_info TWalletInfo) bool {
	er := Db.Create(&wallet_info)
	flag := Db.NewRecord(wallet_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return false
	}
	return true
}

func WalletInfoFind(sql interface{}) (wallet_info_arry []TWalletInfo) {
	Db.Where(sql).Find(&wallet_info_arry)
	return wallet_info_arry
}

func WalletInfoUpdate(wallet_info TWalletInfo, update_info interface{}) bool {
	Db.Model(&wallet_info).Updates(update_info)
	return true
}

func WalletInfoList(wallet_addr, limit_sql string) (list []TWalletInfo) {
	sql := fmt.Sprintf("SELECT * FROM mintklub.t_wallet_info where wallet_addr like '%%%s%%' or user_name like '%%%s%%' %s",
		wallet_addr, wallet_addr, limit_sql)
	Db.Raw(sql).Find(&list)
	return list
}

func WalletInfoLikeCount(wallet_addr string) (total int) {
	sql := fmt.Sprintf("SELECT count(*) FROM mintklub."+
		"t_wallet_info where wallet_addr like '%%%s%%' or user_name like '%%%s%%'",
		wallet_addr, wallet_addr)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("WalletInfoLikeCount error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			logger.Error("WalletInfoLikeCount error", "err", err)
		}
	}
	return total
}
