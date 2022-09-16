package models

import (
	. "ChainServer/models/postgresql"
	"fmt"
	"time"

	"github.com/sea-project/go-logger"
)

type TFavorites struct {
	ID         int
	UserId     int
	NftId      int
	CreateTime time.Time
}

func FavoritesInsert(favorites_info TFavorites) int {
	er := Db.Create(&favorites_info)
	flag := Db.NewRecord(favorites_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return 0
	}
	return favorites_info.ID
}

func FavoritesFind(sql interface{}) (favorites_info_arry []TFavorites) {
	Db.Where(sql).Find(&favorites_info_arry)
	return favorites_info_arry
}

func FavoritesUpdate(favorites_info TFavorites, update_info interface{}) bool {
	Db.Model(&favorites_info).Updates(update_info)
	return true
}

func FavoritesDel(favorites_info TFavorites) {
	Db.Delete(&favorites_info)
}

func FavoritesDelCommon(sql string) {
	Db.Delete(sql)
}

func FavoritesLeftNftLeftMarketList(favorites_addr, params_sql, order_sql, limit_sql string) (market_all_nft []MarketAllNft) {
	sql := `SELECT  (select count(*) from t_favorites as f where n.id = f.nft_id and wallet_addr = '%s') as favorites,
                (SELECT count(1) FROM t_favorites as f where n.id = f.nft_id) as favorites_count, n.nft_name,n.nft_desc,
                n.rights_rules, n.token_id, n.meta_data_uri, n.tx_hash as nft_tx_hash,n.creater as nft_creater,
                n.block_number, n.create_time as nft_create_time, n.media_uri, n.create_tax,n.owner as nft_owner, 
                n.status as nft_status,n.chain_name, n.currency_name, n.lazy, n.media_ipfs_uri, n.collection_id, 
                n.categories_id,m.id as market_id, n.id as nft_id, m.creater as market_creater, m.market_type, 
                m.starting_price, m.end_time,m.reward, m.tx_hash as market_tx_hash,m.cancel_hash, m.deal_hash, 
                m.create_time as market_create_time, m.status as market_status, m.donation,m.donation_addr
             FROM t_favorites as f left join t_nft as n on f.nft_id = n.id left join t_market_list as m on n.id = m.nft_id
                 and (m.status = 1 )
             %s
             %s
             %s
             `
	sql = fmt.Sprintf(sql, favorites_addr, params_sql, order_sql, limit_sql)
	Db.Raw(sql).Find(&market_all_nft)
	return market_all_nft
}

func FavoritesLeftNftLeftMarketListCount(params_sql string) int {
	var total int
	sql := ` SELECT  count(*) FROM t_favorites as f left join t_nft as n on f.nft_id = n.id left join t_market_list 
             as m on n.id = m.nft_id and (m.status = 1 )
             %s`
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
