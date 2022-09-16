package models

import (
	. "ChainServer/models/postgresql"

	"github.com/sea-project/go-logger"
)

type TCategoriesInfo struct {
	Id             int
	CategoriesName string
	CategoriesDesc string
	CollectionId   int
	UserId         int
}

func CategoriesInsert(Categories_info TCategoriesInfo) int {
	er := Db.Create(&Categories_info)
	flag := Db.NewRecord(Categories_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return 0
	}
	return Categories_info.Id
}

func CategoriesFind(sql interface{}) (Categories_info_arry []TCategoriesInfo) {
	Db.Where(sql).Find(&Categories_info_arry)
	return Categories_info_arry
}

func CategoriesUpdate(Categories_info TCategoriesInfo, update_info interface{}) bool {
	Db.Model(&Categories_info).Updates(update_info)
	return true
}
