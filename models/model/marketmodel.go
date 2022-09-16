package model

import (
	"ChainServer/models"
	. "ChainServer/models/postgresql"
	"database/sql"
	"fmt"
	"github.com/sea-project/go-logger"
	"strings"
	"time"
)

// &records.Id, &records.Creater, &records.Tokenid, &records.MarketType, &records.StartingPrice,
// &records.TokenType, &records.EndTime, &records.Bonus, &records.Txhash, &records.CreateTime,
// &records.Sorting, &records.Status
// `id`,`creater` ,`tokenid`,`market_type`,`starting_price`,`token_type`,`end_time`,`bonus`,`txhash`,
//`create_time`,`sorting`,`status`
type MarketListTable struct {
	Id            int    `:"id"`
	SN            string `:"sn"`             // 流水号
	Creater       string `:"creater"`        // 挂单者
	Tokenid       string `:"tokenid"`        // 挂单nft编号
	MarketType    int    `:"market_type"`    // 挂单类型：0 限价购买；1 拍卖
	NftType       int    `:"nft_type"`       // NFT分类
	StartingPrice string `:"starting_price"` // 起拍价
	TokenType     string `:"token_type"`     // 币种：BNB、USDT、CACA
	EndTime       int64  `:"end_time"`       // 拍卖结束时间
	Buyer         string `:"buyer"`          // 购买者
	Bonus         int    `:"bonus"`          // 拍卖参与者分红
	Txhash        string `:"txhash"`         // 交易hash
	CancelHash    string `:"cancel_hash"`    // 取消挂单hash
	DealHash      string `:"deal_hash"`      // deal hash
	CreateTime    int64  `:"create_time"`    // 挂单时间
	Sorting       int    `:"sorting"`        // 排序
	Status        int    `:"status"`         // 挂单状态：0 开启 1 关闭
	ChainType     string `:"chain_type"`     // 链类型
	FreeGas       string `json:"free_gas"`   // 是否免收gas费 0 否 1 是
	Donation      string `json:"donation"`   // 是否捐赠 0 否 1 是
	NftId         int    `json:"nft_id"`     // nft id
}

func MarketInfoFindLimit(sql interface{}, rows, offset int) (token_info_arry []models.TNft) {
	Db.Table("t_nft").Where(sql).Limit(rows).Offset(offset).Find(&token_info_arry)
	return token_info_arry
}

func MarketAllPage(status, mtype string, offset, limit int, nftType, currency, sort, auctionTime, wallet_addr string) (recordList []MarketListTable, total int) {
	var rows *sql.Rows
	var err error
	var querySql string
	if wallet_addr == "" {
		querySql = "SELECT id,sn,creater ,tokenid,market_type, nft_type,starting_price," +
			"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting," +
			"status, chain_type, free_gas, donation,nft_id FROM market_list WHERE"
	} else {
		querySql = "SELECT m.id,m.sn,m.creater ,m.tokenid,m.market_type, m.nft_type,m.starting_price," +
			"m.token_type,m.end_time,m.buyer,m.bonus,m.txhash, m.cancel_hash,m.deal_hash,m.create_time,m.sorting," +
			"m.status, m.chain_type, m.free_gas, m.donation, " +
			"n.id as nft_id FROM market_list as m INNER JOIN t_my_attention as a on m.nft_id = a.nft_id WHERE a." +
			"wallet_addr = '%s' and"
		querySql = fmt.Sprintf(querySql, wallet_addr)
	}

	//whereSql := " sorting >= 0 and market_type = ? and status = ? "
	whereSql := " m.sorting >= 0 "
	nowTime := time.Now().Unix()
	if auctionTime == "1" { // 拍卖中
		whereSql += fmt.Sprintf("and m.end_time > %d ", nowTime)
	} else if auctionTime == "2" { // 拍卖结束
		whereSql += fmt.Sprintf("and m.end_time < %d ", nowTime)
	}
	limitSql := fmt.Sprintf(" LIMIT '%d' OFFSET '%d'", limit, offset)

	orderSql := " order by m.sorting desc "
	if sort == "recently" {
		orderSql = " order by m.create_time desc"
	}
	// ----------------------------------------------------------------
	var mtype_sql string
	if mtype == "" {
		mtype_sql = ""
	} else {
		mtype_sql = "and m.market_type = '%s' "
		mtype_sql = fmt.Sprintf(mtype_sql, mtype)
	}
	whereSql += mtype_sql

	var status_sql string
	if status == "" {
		status_sql = ""
	} else {
		status_sql = "and m.status = '%s' "
		status_sql = fmt.Sprintf(status_sql, status)
	}
	whereSql += status_sql

	if nftType != "" && currency == "" {
		whereSql = whereSql + fmt.Sprintf(" and m.nft_type = '%s' ", nftType)
	} else if nftType == "" && currency != "" {
		whereSql = whereSql + fmt.Sprintf(" and m.token_type = '%s' ", currency)
	} else if nftType != "" && currency != "" {
		whereSql = whereSql + fmt.Sprintf("  and m.nft_type = '%s' and token_type = '%s' ", nftType, currency)
	}
	// 数据
	sql := querySql + whereSql + orderSql + limitSql

	if wallet_addr == "" {
		sql = strings.Replace(sql, "m.", "", -1)
	}
	rows, err = Db.Raw(sql).Rows()

	if err != nil {
		logger.Error("MarketAllPage error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status, &records.ChainType,
			&records.FreeGas, &records.Donation, &records.NftId)
		if err != nil {
			logger.Error("MarketAllPage error", "err", err)
		}
		recordList = append(recordList, records)
	}
	// 总数
	if wallet_addr == "" {
		whereSql = strings.Replace(whereSql, "m.", "", -1)
		rows, err = Db.Raw("SELECT COUNT(*) FROM market_list WHERE " + whereSql).Rows()
	} else {
		start_sql := "SELECT COUNT(*) FROM market_list as m INNER JOIN t_my_attention as a on m.nft_id = a.nft_id WHERE a.wallet_addr = '%s' and "
		start_sql = fmt.Sprintf(start_sql, wallet_addr)
		rows, err = Db.Raw(start_sql + whereSql).Rows()

	}
	if err != nil {
		logger.Error("OrdersCountByMid error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			logger.Error("OrdersCountByMid error", "err", err)
		}
	}
	return recordList, total
}

func MarketByTokenId(nft_id int, status, chain_type string) (recordList []MarketListTable) {
	sql := fmt.Sprintf("SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price,"+
		"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting,"+
		"status FROM market_list WHERE nft_id = %d and status = '%s' and chain_type = '%s'", nft_id, status, chain_type)
	rows, err := Db.Raw(sql).Rows()
	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status)
		if err != nil {
			logger.Error("MarketByTokenId error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func MarketByTokenIdMid(tokenid, mid, status, chain_type string) (recordList []MarketListTable) {
	var rows *sql.Rows
	var err error
	var sql string
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf("and chain_type = '%s'", chain_type)
	}
	if mid == "0" {
		sql = fmt.Sprintf("SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price,"+
			"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting,"+
			"status, free_gas, donation FROM market_list WHERE tokenid = '%s' and status = '%s' ", tokenid, status)
	} else {
		sql = fmt.Sprintf("SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price,"+
			"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting,"+
			"status, free_gas, donation FROM market_list WHERE tokenid = '%s' and id = '%s' and status = '%s' ", tokenid, mid, status)
	}
	sql += chain_type_sql
	rows, err = Db.Raw(sql).Rows()
	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status,
			&records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("MarketByTokenId error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func MarketById(market_id int, chain_type, status string) (recordList []MarketListTable) {
	var rows *sql.Rows
	var err error
	var sql string
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf("and chain_type = '%s'", chain_type)
	}
	sql = fmt.Sprintf("SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price,"+
		"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting,"+
		"status, free_gas, donation, nft_id FROM market_list WHERE id = %d and status = '%s' ", market_id, status)

	sql += chain_type_sql
	rows, err = Db.Raw(sql).Rows()
	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status,
			&records.FreeGas, &records.Donation, &records.NftId)
		if err != nil {
			logger.Error("MarketByTokenId error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// GetMarketHashList 获取有hash、没有更新链信息的记录
func GetMarketHashList(chain_type string) (recordList []MarketListTable) {
	sql := "SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price," +
		"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting," +
		"status FROM market_list WHERE status = 0 and txhash !='' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, chain_type)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("GetMarketHashList error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status)
		if err != nil {
			logger.Error("GetMarketHashList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// GetMarketCancelHashList 获取有取消hash、没有取消挂单的记录
func GetMarketCancelHashList(chain_type, free_gas string) (recordList []MarketListTable) {
	sql := "SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price," +
		"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting," +
		"status FROM market_list WHERE status = 1 and cancel_hash !='' and chain_type = '%s' and free_gas = '%s'"
	sql = fmt.Sprintf(sql, chain_type, free_gas)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("GetMarketCancelHashList error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status)
		if err != nil {
			logger.Error("GetMarketCancelHashList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// GetMarketDealHashList 获取有取消hash、没有取消挂单的记录
func GetMarketDealHashList(chain_type, free_gas string) (recordList []MarketListTable) {
	sql := "SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price," +
		"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting," +
		"status, nft_id FROM market_list WHERE status = 1 and deal_hash !='' and chain_type = '%s' and free_gas = '%s'"
	sql = fmt.Sprintf(sql, chain_type, free_gas)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("GetMarketDealHashList error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status, &records.NftId)
		if err != nil {
			logger.Error("GetMarketDealHashList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// MarketHashUpdate 更新铸造记录hash
func MarketHashUpdate(sn, hash, chain_type string) error {
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf(" and chain_type = '%s'", chain_type)
	}
	sql := "UPDATE market_list SET txhash = '%s' WHERE sn = '%s' and status = 0"
	sql = fmt.Sprintf(sql, hash, sn)
	sql += chain_type_sql
	Db.Exec(sql)
	return nil
}

func MarketHashTokenIdUpdate(id int, token_id string) error {
	sql := "UPDATE market_list SET tokenid = %s WHERE  id = %d"
	sql = fmt.Sprintf(sql, token_id, id)
	Db.Exec(sql)
	return nil
}

// MarketUpdateStatusByTokenid 根据tokenid更新状态
func MarketUpdateStatusByTokenid(tokenid, status, chain_type string) error {
	sql := "UPDATE market_list SET status = '%s' WHERE tokenid = '%s' and status = 1 and chain_type = '%s'"
	sql = fmt.Sprintf(sql, status, tokenid, chain_type)
	Db.Exec(sql)
	return nil
}

func MarketUpdateStatusById(id, txhash string) error {
	sql := "UPDATE market_list SET txhash = '%s' WHERE id = '%s'"
	sql = fmt.Sprintf(sql, txhash, id)
	Db.Exec(sql)
	return nil
}

// MarketCancelHashUpdate 更新取消铸造记录hash
func MarketCancelHashUpdate(tokenid, hash, chain_type string) error {
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf(" and chain_type = '%s'", chain_type)
	}
	sql := "UPDATE market_list SET cancel_hash = '%s' WHERE tokenid = '%s' and status = 1"
	sql = fmt.Sprintf(sql, hash, tokenid)
	sql += chain_type_sql
	Db.Exec(sql)
	return nil
}

// MarketDealHashUpdate 更新deal记录hash
func MarketDealHashUpdate(tokenid, hash, chain_type string) error {
	sql := "UPDATE market_list SET deal_hash = '%s' WHERE tokenid = '%s' and status = 1 and chain_type = '%s'"
	sql = fmt.Sprintf(sql, hash, tokenid, chain_type)
	Db.Exec(sql)
	return nil
}

func MarketDealHashUpdateNew(hash string, market_id int) error {
	sql := "UPDATE market_list SET deal_hash = '%s' WHERE id = %d"
	sql = fmt.Sprintf(sql, hash, market_id)
	Db.Exec(sql)
	return nil
}

// MarketChainUpdate
func MarketChainUpdate(hash, create_time, status, chain_type string) error {
	sql := "UPDATE market_list SET create_time = '%s',status = '%s' " +
		"WHERE txhash = '%s' and status = 0 and chain_type = '%s'"
	sql = fmt.Sprintf(sql, create_time, status, hash, chain_type)
	Db.Exec(sql)
	return nil
}

func MarketStatusUpdateById(id int, status string) error {
	sql := "UPDATE market_list SET status = '%s' " + "WHERE id = %d"
	sql = fmt.Sprintf(sql, status, id)
	Db.Exec(sql)
	return nil
}

// MarketStatusUpdate
func MarketStatusUpdate(hash, status, chain_type string) error {
	sql := "UPDATE market_list SET status = '%s' " + "WHERE txhash = '%s' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, status, hash, chain_type)
	Db.Exec(sql)
	return nil
}

// MarketInsertRecord
func MarketInsertRecord(sn, creater, tokenid, market_type, nft_type, starting_price, token_type, end_time, bonus, sorting, chain_type, free_gas, donation string, status, nft_id, create_time int) error {
	sql := "INSERT INTO market_list(sn,creater,tokenid,market_type,nft_type," +
		"starting_price,token_type,end_time,bonus,sorting,chain_type, free_gas, donation, status, nft_id, create_time) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s', '%s', '%s',%d,%d,%d)"
	sql = fmt.Sprintf(sql, sn, creater, tokenid, market_type, nft_type, starting_price, token_type, end_time, bonus, sorting, chain_type, free_gas, donation, status, nft_id, create_time)
	Db.Exec(sql)
	return nil
}

//  NftGetByTokenId
func MarketGetByMid(mid, chain_type string) (recordList []MarketListTable) {
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf(" and chain_type = '%s'")
	}
	sql := fmt.Sprintf("SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price,"+
		"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting,status,free_gas,donation,nft_id "+
		"FROM market_list WHERE id = '%s'", mid)
	sql += chain_type_sql
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("MarketGetByMid error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status,
			&records.FreeGas, &records.Donation, &records.NftId)
		if err != nil {
			logger.Error("MarketGetByMid error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func MarketHomeLimiteds(market_type, nums, chain_type string) (recordList []MarketListTable) {
	//var chain_type_sql string
	//if chain_type != "" {
	//	chain_type_sql = fmt.Sprintf(" and chain_type = '%s'", chain_type)
	//}

	sql := "SELECT id,creater,market_type,starting_price," +
		"FROM mintklub.t_market_list WHERE market_type = 2 and status = 1"

	//sql = fmt.Sprintf(sql, market_type, nums, chain_type_sql)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("MarketAllPage error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status)
		if err != nil {
			logger.Error("MarketAllPage error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

/*
func MarketList(tokenid, chain_type string) (recordList []MarketListTable) {
	sql := fmt.Sprintf("SELECT id,sn,creater ,tokenid,market_type,nft_type,starting_price,"+
		"token_type,end_time,buyer,bonus,txhash, cancel_hash,deal_hash,create_time,sorting,"+
		"status FROM market_list WHERE tokenid = '%s' and chain_type = '%s'", tokenid, chain_type)
	rows, err := Db.Raw(sql).Rows()
	for rows.Next() {
		var records MarketListTable
		err = rows.Scan(&records.Id, &records.SN, &records.Creater, &records.Tokenid, &records.MarketType, &records.NftType,
			&records.StartingPrice, &records.TokenType, &records.EndTime, &records.Buyer, &records.Bonus,
			&records.Txhash, &records.CancelHash, &records.DealHash, &records.CreateTime, &records.Sorting, &records.Status)
		if err != nil {
			logger.Error("MarketByTokenId error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}*/
//-----------------------------------------------------------------------------------------------------------------------

type MarketList struct {
	Id            int
	SN            string
	Creater       string
	Tokenid       string
	MarketType    int
	NftType       int
	StartingPrice string
	TokenType     string
	EndTime       int64
	Buyer         string
	Bonus         int
	Txhash        string
	CancelHash    string
	DealHash      string
	CreateTime    int64
	Sorting       int
	Status        int
	ChainType     string
	FreeGas       string
	Donation      string
	NftId         int
}

func MarketListCommonFind(page, offset int, data interface{}, order_sql string) (market_list_arry []MarketList) {
	if order_sql == "" {
		Db.Limit(offset).Offset(offset * (page - 1)).Where(data).Find(&market_list_arry)
	} else {
		Db.Order(order_sql).Limit(offset).Offset(offset * (page - 1)).Where(data).Find(&market_list_arry)
	}
	return market_list_arry
}

type MarketAndNft struct {
	Id      int    `json:"id"`
	SN      string `json:"sn"`      // 流水号
	Creater string `json:"creater"` // 挂单者
	//UserName       string `json:"user_name"`       // 创建者用户名
	//ImageUrl       string `json:"image_url"`       // 创建者头像地址
	//OwnerUserName  string `json:"owner_user_name"` // 拥有者用户名
	//OwnerImageUrl  string `json:"owner_image_url"` // 拥有者头像地址
	Tokenid       string `json:"tokenid"`        // 挂单nft编号
	MarketType    int    `json:"market_type"`    // 挂单类型：0 限价购买；1 拍卖
	StartingPrice string `json:"starting_price"` // 起拍价
	TokenType     string `json:"token_type"`     // 币种：BNB、USDT、CACA
	EndTime       string `json:"end_time"`       // 拍卖结束时间
	Buyer         string `json:"buyer"`          // 拍卖参与者分红
	Bonus         int    `json:"bonus"`          // 拍卖参与者分红
	Txhash        string `json:"txhash"`         // 交易hash
	CancelHash    string `json:"cancel_hash"`    // 取消挂单hash
	DealHash      string `json:"deal_hash"`      // 取消挂单hash
	CreateTime    string `json:"create_time"`    // 挂单时间
	Sorting       int    `json:"sorting"`        // 排序
	Status        int    `json:"status"`         // 挂单状态：0 开启 1 关闭

	NftName        string `json:"nft_name"`        // 名称
	NftDesc        string `json:"nft_desc"`        // 描述
	TokenId        string `json:"token_id"`        // nft编号
	TokenUri       string `json:"token_uri"`       // nft地址
	NftTxhash      string `json:"nft_txhash"`      // 铸造nft的hash
	NftCreater     string `json:"nft_creater"`     // nft创建者
	CreateNumber   int    `json:"create_number"`   // 创建区块高度
	Nft_CreateTime string `json:"nft_create_time"` // 创建时间戳
	MediaUri       string `json:"media_uri"`       // 数字资产地址
	CreateTax      int    `json:"create_tax"`      // 铸造税
	Owner          string `json:"owner"`           // nft拥有者
	NftType        int    `json:"nft_type"`        // nft分类
	Approved       int    `json:"approved"`        // nft审核状态
	//OfferCount     int    `json:"offer_count"`     // 竞拍出价次数
	//Attention      int    `json:"attention"`       // 是否关注
	//AttentionCount int    `json:"attention_count"` // 关注次数
	ChainType string `json:"chain_type"` // 链类型
	FreeGas   string `json:"free_gas"`   // 是否免收gas费 0 否 1 是
	Donation  string `json:"donation"`   // 是否捐赠 0 否 1 是
	NftId     int    `json:"nft_id"`     // nft id
}

type MarketAndNftAll struct {
	MarketAndNft
	UserName      string `json:"user_name"`       // 创建者用户名
	ImageUrl      string `json:"image_url"`       // 创建者头像地址
	OwnerUserName string `json:"owner_user_name"` // 拥有者用户名
	OwnerImageUrl string `json:"owner_image_url"` // 拥有者头像地址
}

func Test(screen_sql, order_sql, limit_sql string) *[]MarketAndNftAll {
	sql := `SELECT m.id, m.sn, m.creater, m.tokenid, m.market_type, m.nft_type, m.starting_price, m.token_type, m.end_time,
                m.buyer, m.bonus, m.txhash, m.cancel_hash,m.deal_hash, m.create_time, m.sorting, m.status,
                n.nft_name, n.nft_desc,n.token_id,n.token_uri,n.txhash,n.creater,n.create_number,n.create_time,
                n.media_uri,n.create_tax,n.owner,n.nft_type,n.approved,n.chain_type,n.free_gas,n.donation,n.id
            FROM market_list as m inner join nft as n on m.nft_id = n.id `
	sql += screen_sql
	sql += limit_sql
	sql += order_sql
	var market_nft []MarketAndNftAll
	data := Db.Raw(sql).Find(&market_nft)
	return data.Value.(*[]MarketAndNftAll)
}

func TestRight(screen_sql, order_sql, limit_sql string) *[]MarketAndNftAll {
	sql := `SELECT m.id, m.sn, m.creater, m.tokenid, m.market_type, m.starting_price, m.token_type, m.end_time,
                m.buyer, m.bonus, m.txhash, m.cancel_hash,m.deal_hash, m.create_time, m.sorting, m.status,
                n.nft_name, n.nft_desc,n.token_id,n.token_uri,n.txhash,n.creater,n.create_number,n.create_time,
                n.media_uri,n.create_tax,n.owner,n.nft_type,n.approved,n.chain_type,n.free_gas,n.donation,n.id
            FROM market_list as m right join nft as n on m.nft_id = n.id `
	sql += screen_sql
	sql += limit_sql
	sql += order_sql
	var market_nft []MarketAndNftAll
	data := Db.Raw(sql).Find(&market_nft)
	return data.Value.(*[]MarketAndNftAll)
}

/*
SELECT m.id, m.sn, m.creater, m.tokenid, m.market_type, m.starting_price, m.token_type, m.end_time, m.buyer, m.bonus, m.txhash, m.cancel_hash, m.create_time, m.sorting, m.status,
              n.nft_name, n.nft_desc,n.token_id,n.token_uri,n.nft_txhash,n.nft_creater,n.create_number,n.nft_create_time,n.media_uri,n.create_tax,n.owner,n.nft_type,n.approved,n.offer_count,n.attention,
  n.attention_count,n.chain_type,n.free_gas,n.donation,n.id
			FROM market_list as m inner join nft as n on m.nft_id = n.id
*/
