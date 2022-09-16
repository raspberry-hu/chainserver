package models

import (
	. "ChainServer/models/postgresql"
	"github.com/sea-project/go-logger"
)

type TPermissions struct {
	Id         int
	WalletAddr string
}

func TPermissionsInsert(t_permissions TPermissions) int {
	er := Db.Create(&t_permissions)
	flag := Db.NewRecord(t_permissions) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return 0
	}
	return t_permissions.Id
}

func TPermissionsFind(sql interface{}) (t_permissions_arry []TPermissions) {
	Db.Where(sql).Find(&t_permissions_arry)
	return t_permissions_arry
}
