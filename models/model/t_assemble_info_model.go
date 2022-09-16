package model

import (
	. "ChainServer/models/postgresql"
	"github.com/sea-project/go-logger"
)

type TAssembleInfo struct {
	Id           int64
	WalletAddr   string
	AssembleName string
}

func TAssembleInfoInsert(t_assemble_info TAssembleInfo) bool {
	er := Db.Create(&t_assemble_info)
	flag := Db.NewRecord(t_assemble_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return false
	}
	return true
}

func TemplateFind(data map[string]interface{}) (t_assemble_info_arry []TAssembleInfo) {
	Db.Order("r_time desc").Where(data).Find(&t_assemble_info_arry)
	return t_assemble_info_arry
}

func TAssembleInfoUpdate(t_assemble_info TAssembleInfo, data map[string]interface{}) bool {
	Db.Model(&t_assemble_info).Updates(data)
	return true
}

/*
func TemplateDel(t_assemble_info TAssembleInfo) bool {
	Db.Where(template).Find(&template).Delete(template)
	return true
}
*/
