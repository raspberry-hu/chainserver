package model

import (
	. "ChainServer/models/postgresql"
	"fmt"
	"github.com/sea-project/go-logger"
)

type Banner struct {
	Id         int    `json:"id"`
	BannerName string `json:"banner_name"` // banner名称
	BannerImg  string `json:"banner_img"`  // banner图片地址
	BannerUrl  string `json:"banner_url"`  // banner链接地址
}

// BannerList 所有banner
func BannerList(chain_type string) (recordList []Banner) {
	// rows, err := Db.Query("SELECT id,banner_name,banner_img ,banner_url FROM caca.banner")
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf(" and chain_type = '%s'", chain_type)
	}
	sql := ""
	sql += chain_type_sql
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("BannerList error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records Banner
		err = rows.Scan(&records.Id, &records.BannerName, &records.BannerImg, &records.BannerUrl)
		if err != nil {
			logger.Error("BannerList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}
