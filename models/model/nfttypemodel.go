package model

import (
	. "ChainServer/models/postgresql"
	"fmt"
	"github.com/sea-project/go-logger"
)

type NftType struct {
	Id       int    `json:"id"`
	TypeName string `json:"type_name"` // 分类名称
}

// TypeList 所有nft type
func TypeList(chain_type string) (recordList []NftType) {
	//rows, err := DB.Query("SELECT id,type_name FROM nft_type")
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf(" and chain_type = '%s'", chain_type)
	}
	sql := "SELECT id,type_name FROM nft_type"
	sql += chain_type_sql
	rows, err := Db.Raw(sql).Rows()

	if err != nil {
		logger.Error("TypeList error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records NftType
		err = rows.Scan(&records.Id, &records.TypeName)
		if err != nil {
			logger.Error("TypeList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}
