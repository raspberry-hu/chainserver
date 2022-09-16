package models

import (
	. "ChainServer/models/postgresql"
	"database/sql"
	"fmt"

	"github.com/sea-project/go-logger"
)

type TOrder struct {
	Id           int    `json:"id"`
	NftId        int    `json:"nft_id"`
	MarketId     int    `json:"market_id"`
	MarketType   int    `json:"market_type"`
	TokenId      string `json:"token_id"`
	CreateTime   int    `json:"create_time"`
	OnchainTime  int    `json:"onchain_time"`
	TxHash       string `json:"tx_hash"`
	Buyer        int    `json:"buyer"`
	Seller       int    `json:"seller"`
	Price        string `json:"price"`
	ChainName    string `json:"chain_name"`
	CurrencyName string `json:"currency_name"`
	Lazy         int    `json:"lazy"`
	Status       int    `json:"status"`
	Donation     int    `json:"donation"`
}

func OrderInsert(order_info TOrder) int {
	er := Db.Create(&order_info)
	flag := Db.NewRecord(order_info) // 主键是否为空
	if er.Error != nil || flag {
		logger.Error("insert xy_template err failed", er.Error)
		return 0
	}
	return order_info.Id
}

func OrderFind(sql interface{}) (order_info_arry []TOrder) {
	Db.Where(sql).Find(&order_info_arry)
	return order_info_arry
}

func OrderFindCount(sql interface{}) (count int) {
	Db.Model(&TOrder{}).Where(sql).Count(&count)
	return count
}

//
func OrderFindLimit(sql interface{}, rows, offset int) (order_info_arry []TOrder) {
	Db.Order("create_time desc").Limit(rows).Offset(offset).Where(sql).Find(&order_info_arry)
	return order_info_arry
}

func OrderUpdate(order_info TOrder, update_info interface{}) bool {
	Db.Model(&order_info).Updates(update_info)
	return true
}

func OrdersGetMaxPrice(market_id string) (buyer string, price float32) {
	var rows *sql.Rows
	var err error
	sql := fmt.Sprintf("SELECT buyer,SUM(cast (price as float8)) AS total FROM mintklub.t_order WHERE market_id = %s "+
		" and status = 1 GROUP BY buyer ORDER BY total DESC LIMIT 1", market_id)
	rows, err = Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("OrdersGetMaxPrice error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&buyer, &price)

		if err != nil {
			logger.Error("OrdersGetMaxPrice error", "err", err)
		}
	}
	return buyer, price
}
