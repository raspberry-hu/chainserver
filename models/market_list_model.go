package models

import (
	. "ChainServer/models/postgresql"
	"fmt"

	"github.com/sea-project/go-logger"
)

type TMarketList struct {
	Id             int    `json:"id"`
	NftId          int    `json:"nft_id"`
	Creater        int    `json:"creater"`
	TokenId        string `json:"token_id"`
	MarketType     int    `json:"market_type"`
	StartingPrice  string `json:"starting_price"`
	EndTime        int64  `json:"end_time"`
	Buyer          int    `json:"buyer"`
	Reward         int    `json:"reward"`
	TxHash         string `json:"tx_hash"`
	CancelHash     string `json:"cancel_hash"`
	DealHash       string `json:"deal_hash"`
	CreateTime     int64  `json:"create_time"`
	Status         int    `json:"status"`
	ChainName      string `json:"chain_name"`
	CurrencyName   string `json:"currency_name"`
	Lazy           int    `json:"lazy"` // 1
	Donation       int    `json:"donation"`
	DonationUserId int    `json:"donation_userid"`
}

func MarketListInsert(market_info TMarketList) int {
	er := Db.Create(&market_info)
	flag := Db.NewRecord(market_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return 0
	}
	return market_info.Id
}

func MarketListFind(sql interface{}) (market_info_arry []TMarketList) {
	Db.Where(sql).Find(&market_info_arry)
	return market_info_arry
}

func MarketListFindLimit(page, offset int, sql interface{}) (market_info_arry []TMarketList) {
	Db.Limit(page).Offset(offset).Where(sql).Find(&market_info_arry)
	return market_info_arry
}

func MarketListFindCount(sql interface{}) (count int) {
	Db.Model(&TMarketList{}).Where(sql).Count(&count)
	return count
}

func MarketCollectionListFind(nft_id int, creator int, min, max string) (market_info_arry []TMarketList) {
	sql := `SELECT * FROM mintklub.t_market_list WHERE %s`
	var params_sql string
	if min == "" && max != "" {
		params_sql = fmt.Sprintf("nft_id=%d AND creater='%d' AND starting_price<='%s'", nft_id, creator, max)
	} else if min != "" && max == "" {
		params_sql = fmt.Sprintf("nft_id=%d AND creater='%d' AND starting_price>='%s'", nft_id, creator, min)
	} else if min != "" && max != "" {
		params_sql = fmt.Sprintf("nft_id=%d AND creater='%d' AND starting_price>='%s' AND starting_price<='%s'", nft_id, creator, min, max)
	} else {
		params_sql = fmt.Sprintf("nft_id=%d AND creater='%d'", nft_id, creator)
	}
	sql = fmt.Sprintf(sql, params_sql)
	Db.Raw(sql).Find(&market_info_arry)
	return market_info_arry
}

func MarketListUpdate(market_info TMarketList, update_info interface{}) bool {
	Db.Model(&market_info).Updates(update_info)
	return true
}

type MarketAllNft struct {
	FavoritesCount string `json:"favorites_count"`
	Favorites      string `json:"favorites"`
	NftName        string `json:"nft_name"`
	NftDesc        string `json:"nft_desc"`
	RightsRules    string `json:"rights_rules"`
	TokenId        string `json:"token_id"`
	MetaDataUri    string `json:"meta_data_uri"`
	NftTxHash      string `json:"nft_tx_hash"`
	NftCreater     int    `json:"nft_creater"`
	BlockNumber    int    `json:"block_number"`
	NftCreateTime  int    `json:"nft_create_time"`
	MediaUri       string `json:"media_uri"`
	CreateTax      string `json:"create_tax"`
	NftOwner       int    `json:"nft_owner"`
	NftStatus      int    `json:"nft_status"`
	ChainName      string `json:"chain_name"`
	CurrencyName   string `json:"currency_name"`
	Lazy           int    `json:"lazy"`
	CollectionId   int    `json:"collection_id"`
	CategoriesId   int    `json:"categories_id"`
	MediaIpfsUri   string `json:"media_ipfs_uri"`
	ExploreUri     string `json:"explore_uri"` // 音视频的额外展示uri

	MarketId         int    `json:"market_id"`
	NftId            int    `json:"nft_id"`
	MarketCreater    string `json:"market_creater"`
	MarketType       int    `json:"market_type"`
	StartingPrice    string `json:"starting_price"`
	EndTime          int    `json:"end_time"`
	Reward           int    `json:"reward"`
	MarketTxHash     string `json:"market_tx_hash"`
	CancelHash       string `json:"cancel_hash"`
	DealHash         string `json:"deal_hash"`
	MarketCreateTime int    `json:"market_create_time"`
	MarketStatus     int    `json:"market_status"`
	Donation         int    `json:"donation"`
	DonationUserId   int    `json:"donation_userid"`
}

func MarketAllFind(user_id int, params_sql, order_sql, limit_sql string) (market_all_nft []MarketAllNft) {
	sql := `SELECT  (select count(*) from mintklub.t_favorites as f where n.id = f.nft_id and user_id = '%d') as favorites,
(SELECT count(1) FROM mintklub.t_favorites as f where n.id = f.nft_id) as favorites_count, n.nft_name, 
n.nft_desc, n.rights_rules, n.token_id, n.meta_data_uri, n.tx_hash as nft_tx_hash,
n.creater as nft_creater, n.block_number, n.create_time as nft_create_time, n.media_uri, n.create_tax, n.owner as nft_owner, n.status as nft_status,
n.chain_name, n.currency_name, n.lazy, n.media_ipfs_uri, n.explore_uri, n.collection_id, n.categories_id,
m.id as market_id, n.id as nft_id, m.creater as market_creater, m.market_type, m.starting_price, m.end_time, m.reward, 
m.tx_hash as market_tx_hash,
m.cancel_hash, m.deal_hash, m.create_time as market_create_time, m.status as market_status, m.donation, m.donation_user_id 
FROM mintklub.t_nft as n inner join mintklub.t_market_list as m on m.nft_id = n.id %s %s %s
`
	sql = fmt.Sprintf(sql, user_id, params_sql, order_sql, limit_sql)
	Db.Raw(sql).Find(&market_all_nft)
	return market_all_nft
}

func MarketALeftNftFind(user_id int, params_sql, order_sql, limit_sql string) (market_all_nft []MarketAllNft) {
	sql := `SELECT  (select count(*) from mintklub.t_favorites as f where n.id = f.nft_id and user_id = '%d') as favorites,
            (SELECT count(1) FROM mintklub.t_favorites as f where n.id = f.nft_id) as favorites_count, n.nft_name, 
            n.nft_desc, n.rights_rules, n.token_id, n.meta_data_uri, n.tx_hash as nft_tx_hash,
            n.creater as nft_creater, n.block_number, n.create_time as nft_create_time, n.media_uri, n.create_tax, 
            n.owner as nft_owner, n.status as nft_status,n.chain_name, n.currency_name, n.lazy, n.media_ipfs_uri, 
            n.collection_id, n.categories_id, n.explore_uri, m.id as market_id, n.id as nft_id, m.creater as market_creater, 
m.market_type, 
            m.starting_price, m.end_time,m.reward, m.tx_hash as market_tx_hash,m.cancel_hash, m.deal_hash, 
            m.create_time as market_create_time, m.status as market_status, m.donation, m.donation_user_id 
            FROM mintklub.t_nft as n left join t_market_list as m on m.nft_id = n.id and (m.status = 1 )
            %s
            %s
            %s
`
	sql = fmt.Sprintf(sql, user_id, params_sql, order_sql, limit_sql)
	Db.Raw(sql).Find(&market_all_nft)
	return market_all_nft
}

func MarketAllCount(params_sql string) int {
	var total int
	sql := ` SELECT  count(*) FROM mintklub.t_nft as n inner join t_market_list as m on m.nft_id = n.id %s`
	sql = fmt.Sprintf(sql, params_sql)
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

func NftLeftMarketCount(params_sql string) int {
	var total int
	sql := ` SELECT  count(*) FROM mintklub.t_nft as n left join t_market_list as m on m.nft_id = n.id and (m.
status = 2  ) %s`
	sql = fmt.Sprintf(sql, params_sql)
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
