package models

import (
	. "ChainServer/models/postgresql"
	"fmt"

	"github.com/sea-project/go-logger"
)

type TCollectionInfo struct {
	Id               int     `json:"id"`
	UserId           int     `json:"user_id"`
	CollectionName   string  `json:"collection_name"`
	ChainName        string  `json:"chain_name"`
	CurrencyName     string  `json:"currency_name"`
	LogoImage        string  `json:"logo_image"`
	FeaturedImageUrl string  `json:"featured_image_url"`
	BannerImageUrl   string  `json:"banner_image_url"`
	CollectionDesc   string  `json:"collection_desc"`
	CreateTax        int     `json:"create_tax"`
	Category         string  `json:"category"`
	Items            int     `json:"items"`
	Favorites        int     `json:"favorites"`
	Amount           float32 `json:"amount"`
	TxHash           string  `json:"tx_hash"` //1
	Status           int     `json:"status"`  //1
}

func CollectionInsert(collection_info TCollectionInfo) int {
	er := Db.Create(&collection_info)
	flag := Db.NewRecord(collection_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return 0
	}
	return collection_info.Id
}

func CollectionFind(sql interface{}) (collection_info_arry []TCollectionInfo) {
	Db.Table("ant.t_collection_info").Where(sql).Find(&collection_info_arry)
	return collection_info_arry
}

func CollectionUpdate(collection_info TCollectionInfo, update_info interface{}) bool {
	Db.Model(&collection_info).Updates(update_info)
	return true
}

func UpdateItemsByCollection(collectionInfo TCollectionInfo, updateInfo interface{}) bool {
	//	每次铸币更新
	Db.Model(&collectionInfo).Update(updateInfo)
	return true
}

func UpdateFavoriteByCollection(collectionInfo TCollectionInfo, updateInfo interface{}) bool {
	//	每次收藏更新
	Db.Model(&collectionInfo).Update(updateInfo)
	return true
}

func UpdateAmountByCollection(collectionInfo TCollectionInfo, updateInfo interface{}) bool {
	//	每次创建订单更新
	Db.Model(&collectionInfo).Update(updateInfo)
	return true
}

func CollectionListFind(contains_sql, limit_sql string) (collection_list_arry []CollectionList) {
	sql := `
SELECT
	logo_image,
	collection_name,
	collection_desc,
	banner_image_url,
	( SELECT COUNT ( * ) FROM t_nft AS n WHERE n.collection_id = C.ID ) AS nft_count,
	user_id,
	( SELECT user_name FROM t_user_info AS T WHERE C.user_id = T.id ) AS user_name 
FROM
	t_collection_info AS C 
    %s
	%s
`
	sql = fmt.Sprintf(sql, contains_sql, limit_sql)
	Db.Raw(sql).Find(&collection_list_arry)
	return collection_list_arry
}

func CollectionCount(contains_sql string) (total int) {
	sql := `
              SELECT count(*) FROM t_collection_info as C 
              %s
           `
	sql = fmt.Sprintf(sql, contains_sql)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("MarketAllCount error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			logger.Error("MarketAllCount error", "err", err)
		}
	}
	return total
}

type CollectionList struct {
	LogoImage        string `json:"logo_image"`
	CollectionName   string `json:"collection_name"`
	CollectionDesc   string `json:"collection_desc"`
	FeaturedImageUrl string `json:"featured_image_url"`
	NftCount         int    `json:"nft_count"`
	UserId           int    `json:"user_id"`
	UserName         string `json:"user_name"`
}

type WalletAddrList struct {
	UserId string `json:"user_id"`
}

type UserIdList struct {
	UserId int `json:"user_id"`
}
