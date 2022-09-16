package models

import (
	. "ChainServer/models/postgresql"
	"fmt"

	"github.com/sea-project/go-logger"
)

type TCollectRanking struct {
	Id             int    `json:"id"`
	CollectionId   int    `json:"collection_id"`
	LogoImageURL   string `json:"logo_image_url"`
	BannerImageURL string `json:"banner_image_url"`
	CollectionName string `json:"collection_name"`
	CollectionDesc string `json:"collection_desc"`
	CreateTax      int    `json:"create_tax"`
	CategoryName   string `json:"category_name"`
	ChainName      string `json:"chain_name"`
	Owner          string `json:"owner"`
	Items          int    `json:"items"`
	Favorites      int    `json:"favorites"`
}

// 插入排列数据
func RankingInsert(rank_info TCollectRanking) int {
	err := Db.Create(&rank_info)
	// 查看主键是否为空
	flag := Db.NewRecord(rank_info)
	if err.Error != nil || flag {
		logger.Error("insert collection_ranking failed", err.Error)
		return 0
	}
	return rank_info.Id
}

// 返回所有元组
func RankingFindAll(limit string) (rankInfoArray []TCollectionInfo) {

	sql := `SELECT * FROM ant.t_collection_info %s`
	sql = fmt.Sprintf(sql, limit)
	result := Db.Raw(sql).Find(&rankInfoArray)
	fmt.Println("数据库中的所有数据：")
	fmt.Println(rankInfoArray)
	// 返回数据的行数
	affected := result.RowsAffected
	fmt.Printf("\n数据库中的所有数据：%d\n", affected)
	return rankInfoArray
}

// 返回指定的元组
func RankingByCondition(category, chain string, limit string) (rankInfoArray []TCollectionInfo) {
	var condition string

	if category != "" && chain != "" {
		condition = fmt.Sprintf("WHERE C.category='%s' AND C.chain_name='%s' %s", category, chain, limit)
	} else if category == "" && chain != "" {
		condition = fmt.Sprintf("WHERE C.chain_name='%s' %s", chain, limit)
	} else if category != "" && chain == "" {
		condition = fmt.Sprintf("WHERE C.category='%s' %s", category, limit)
	} else {
		condition = limit
	}

	sql := `SELECT * FROM ant.t_collection_info AS C %s`
	sql = fmt.Sprintf(sql, condition)
	Db.Raw(sql).Find(&rankInfoArray)
	return rankInfoArray
}
