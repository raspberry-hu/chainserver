package models

import (
	. "ChainServer/models/postgresql"
	"fmt"

	"github.com/sea-project/go-logger"
)

type TNft struct {
	Id           int    `json:"id"`
	Sn           string `json:"sn"`
	NftName      string `json:"nft_name"`
	NftDesc      string `json:"nft_desc"`
	RightsRules  string `json:"rights_rules"`
	TokenId      string `json:"token_id"`
	MetaDataUri  string `json:"meta_data_uri"`
	TxHash       string `json:"tx_hash"`
	TransferHash string `json:"transfer_hash"` // 蚂蚁链中交易返回的hash值
	Creater      int    `json:"creater"`
	BlockNumber  int    `json:"block_number"`
	CreateTime   int64  `json:"create_time"` // 铸造时间
	BuyTime      int    `json:"buy_time"`    // 购买时间
	MediaUri     string `json:"media_uri"`
	CreateTax    int    `json:"create_tax"`
	Owner        int    `json:"owner"`
	Status       int    `json:"status"`
	MarketType   int    `json:"market_type"`
	Approved     int    `json:"approved"`
	ChainName    string `json:"chain_name"`
	CurrencyName string `json:"currency_name"`
	Lazy         int64  `json:"lazy"` //1
	MediaIpfsUri string `json:"media_ipfs_uri"`
	CollectionId int    `json:"collection_id"`
	CategoriesId int    `json:"categories_id"`
	ExploreUri   string `json:"explore_uri"`
	AntNFTBuyer  string `json:"ant_nft_buyer"` // 蚂蚁链nft购买者
	AntTokenId   int    `json:"ant_token_id"`  // 蚂蚁链的tokenId
	AntCount     int    `json:"ant_count"`     // 蚂蚁链的铸造个数
	AntNftOwner  string `json:"ant_nft_owner"` // 蚂蚁链的nft铸造者
	AntTokenUrl  string `json:"ant_token_url"` // 蚂蚁链的tokenUrl
	AntTxHash    string `json:"ant_tx_hash"`   // 蚂蚁链藏品对应公链的hash值
	AntNftId     int    `json:"ant_nft_id"`    // 蚂蚁链一个藏品种类的序列值
	AntPrice     int    `json:"ant_price"`     // 蚂蚁链藏品价格
}

func NftInsert(nft_info TNft) int {
	er := Db.Table("ant.t_nft").Create(&nft_info)
	//er = Db.Table("ant.t_nft")
	flag := Db.NewRecord(nft_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return 0
	}
	return nft_info.Id
}

func NftFind(sql interface{}) (nft_info_arry []TNft) {
	Db.Table("ant.t_nft").Where(sql).Find(&nft_info_arry)
	return nft_info_arry
}

func NftFindByLimit(sql interface{}, rows, offset int) (nft_info_arry []TNft) {
	Db.Table("ant.t_nft").Where(sql).Limit(rows).Offset(offset).Find(&nft_info_arry)
	return nft_info_arry
}

func NftUpdate(nft_info TNft, update_info interface{}) bool {
	Db.Model(&nft_info).Updates(update_info)
	return true
}

func CollectionNftFind(collection_id int, sell_type int, currency_name string, name string, owner int) (nft_info_arry []TNft) {
	sql := `SELECT * FROM ant.t_nft %s`
	var params_sql string
	if sell_type == 0 {
		if currency_name != "" && name != "" {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND owner='%d' AND currency_name='%s' AND nft_name LIKE '%%%s%%'", collection_id, owner, currency_name, name)
		} else if currency_name != "" && name == "" {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND owner='%d' AND currency_name='%s'", collection_id, owner, currency_name)
		} else if currency_name == "" && name != "" {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND owner='%d' AND nft_name LIKE '%%%s%%'", collection_id, owner, name)
		} else {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND owner='%d'", collection_id, owner)
		}
	} else {
		if currency_name != "" && name != "" {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND market_type=%d AND currency_name='%s' AND nft_name LIKE '%%%s%%'", collection_id, sell_type, currency_name, name)
		} else if currency_name != "" && name == "" {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND market_type=%d AND currency_name='%s'", collection_id, sell_type, currency_name)
		} else if currency_name == "" && name != "" {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND market_type=%d AND nft_name LIKE '%%%s%%'", collection_id, sell_type, name)
		} else {
			params_sql = fmt.Sprintf("WHERE collection_id=%d AND market_type=%d", collection_id, sell_type)
		}
	}

	sql = fmt.Sprintf(sql, params_sql)
	Db.Raw(sql).Find(&nft_info_arry)
	return nft_info_arry
}
