package model

import (
	. "ChainServer/models/postgresql"
	"database/sql"
	"fmt"
	"github.com/sea-project/go-logger"
	"time"
)

// &records.Id, &records.Mid, &records.Tokenid, &records.CreateTime, &records.OnchainTime,
// &records.Txhash, &records.Buyer, &records.Seller, &records.Price, &records.TokenType,
// &records.Status
// `id`,`mid`,`tokenid`,`create_time`,`onchain_time`,`txhash`,`buyer`,`seller`,`price`,`token_type`,`status`
type OrderTable struct {
	Id          int    `:"id"`
	SN          string `:"sn"`           // 流水号
	Mid         int    `:"mid"`          // 对应的挂单id
	MType       int    `:"mtype"`        // 对应的挂单类型
	Tokenid     string `:"tokenid"`      // tokenid
	CreateTime  int64  `:"create_time"`  // 下单时间
	OnchainTime int64  `:"onchain_time"` // 上链时间
	Txhash      string `:"txhash"`       // 交易hash
	Buyer       string `:"buyer"`        // 买家
	Seller      string `:"seller"`       // 卖家
	Price       string `:"price"`        // 价格
	TokenType   string `:"token_type"`   // 币种：BNB、USDT、CACA
	Status      int    `:"status"`       // 订单状态：0 已提交 1 已上链 2 失败
	ChainType   string `:"chain_type"`   // 链类型
	FreeGas     string `json:"free_gas"` // 是否免收gas费 0 否 1 是
	Donation    string `json:"donation"` // 是否捐赠 0 否 1 是
	NftId       int    `json:"nft_id"`   // 是否捐赠 0 否 1 是
}

// OrderHashUpdate 更新铸造记录hash
func OrderHashUpdate(sn, hash, chain_type string) error {
	// 异常捕获
	//defer func() {
	//	if r := recover(); r != nil {
	//		logger.Error("OrderHashUpdate error", "err", r)
	//	}
	//}()
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf(" and chain_type = '%s'", chain_type)
	}
	sql := "UPDATE t_order SET txhash = '%s' WHERE sn = '%s' and status = 0"
	sql = fmt.Sprintf(sql, hash, sn)
	sql += chain_type_sql
	Db.Exec(sql)

	return nil
}

func OrderHashUpdateById(hash, order_id string) {
	sql := "UPDATE t_order SET txhash = '%s' WHERE id = %s and status = 0"
	sql = fmt.Sprintf(sql, hash, order_id)
	sql = fmt.Sprintf(sql)
	Db.Exec(sql)
}

// OrderChainUpdate
func OrderChainUpdate(hash, create_time, status, chain_type string) error {
	// 异常捕获
	sql := "UPDATE t_order SET onchain_time = '%s',status = '%s' " +
		"WHERE txhash = '%s' and status = '0' and chain_type = '%s'"
	// 执行修改
	sql = fmt.Sprintf(sql, create_time, status, hash, chain_type)
	Db.Exec(sql)
	return nil
}

func OrderChainUpdateTokenId(tokenid string, id int) {
	// 异常捕获
	sql := "UPDATE t_order SET tokenid = '%s' WHERE id = %d"
	// 执行修改
	sql = fmt.Sprintf(sql, tokenid, id)
	Db.Exec(sql)
}

// InsertOrderChainInfo
func InsertOrderChainInfo(hash, onchain_time, status string) error {
	// 异常捕获
	defer func() {
		if r := recover(); r != nil {
			logger.Error("InsertOrderChainInfo error", "err", r)
		}
	}()
	sql := "UPDATE t_order SET onchain_time = '%s',status = '%s' " +
		"WHERE txhash = '%s' and status = 0"
	sql = fmt.Sprintf(sql, onchain_time, status, hash)
	// 执行修改
	Db.Exec(sql)
	return nil
}

// GetOrderHashList 获取有hash、没有更新链信息的记录
func GetOrderHashList(chain_type string) (recordList []OrderTable) {

	sql := "SELECT id,sn,mid,mtype,tokenid,create_time,onchain_time,txhash," +
		"buyer,seller,price,token_type,status FROM t_order WHERE status = 0 and txhash !='' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, chain_type)

	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("GetOrderHashList error", "err", err)
	}
	for rows.Next() {
		var records OrderTable
		err = rows.Scan(&records.Id, &records.SN, &records.Mid, &records.MType, &records.Tokenid, &records.CreateTime,
			&records.OnchainTime, &records.Txhash, &records.Buyer, &records.Seller, &records.Price,
			&records.TokenType, &records.Status)
		if err != nil {
			logger.Error("GetOrderHashList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func GetOrderHashListNew(chain_type, free_gas string) (recordList []OrderTable) {

	sql := "SELECT id,sn,mid,mtype,tokenid,create_time,onchain_time,txhash," +
		"buyer,seller,price,token_type,status,free_gas,donation, nft_id FROM t_order WHERE status = 0 and txhash !='' and chain_type = '%s' and free_gas = '%s'"
	sql = fmt.Sprintf(sql, chain_type, free_gas)

	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("GetOrderHashList error", "err", err)
	}
	for rows.Next() {
		var records OrderTable
		err = rows.Scan(&records.Id, &records.SN, &records.Mid, &records.MType, &records.Tokenid, &records.CreateTime,
			&records.OnchainTime, &records.Txhash, &records.Buyer, &records.Seller, &records.Price,
			&records.TokenType, &records.Status, &records.FreeGas, &records.Donation, &records.NftId)
		if err != nil {
			logger.Error("GetOrderHashList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// OrderInsertRecord
func OrderInsertRecord(sn, mid, mtype, tokenid, buyer, seller, price, token_type, chain_type, free_gas, donation string, status, nft_id int) error {
	create_time := time.Now().Unix()
	//create_time := common.IntToDate(time.Now().Unix())
	sql := "INSERT INTO t_order(sn,mid,mtype,tokenid,create_time,buyer,seller,price,token_type, chain_type,free_gas, donation,status, nft_id) VALUES('%s','%s','%s','%s','%d','%s','%s','%s','%s','%s', '%s','%s',%d, %d)"
	sql = fmt.Sprintf(sql, sn, mid, mtype, tokenid, create_time, buyer, seller, price, token_type, chain_type, free_gas, donation, status, nft_id)
	Db.Exec(sql)
	return nil
}

//  NftGetByTokenId
func OrderGetByOid(oid, chain_type string) (recordList []OrderTable) {
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf(" and chain_type = '%s'", chain_type)
	}

	sql := fmt.Sprintf("SELECT id,sn,mid,mtype,tokenid,create_time,onchain_time,txhash,"+
		"buyer,seller,price,token_type,status, free_gas, donation FROM t_order WHERE id = '%s'", oid)
	sql += chain_type_sql
	rows, err := Db.Raw(sql).Rows()
	for rows.Next() {
		var records OrderTable
		err = rows.Scan(&records.Id, &records.SN, &records.Mid, &records.MType, &records.Tokenid, &records.CreateTime,
			&records.OnchainTime, &records.Txhash, &records.Buyer, &records.Seller, &records.Price,
			&records.TokenType, &records.Status, &records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("OrderGetByOid error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func OrderGetByMidTokenId(mid string) (recordList []OrderTable) {

	sql := fmt.Sprintf("SELECT id,sn,mid,mtype,tokenid,create_time,onchain_time,txhash,"+
		"buyer,seller,price,token_type,status, free_gas, donation FROM t_order WHERE mid = '%s' and tokenid != ''", mid)
	rows, err := Db.Raw(sql).Rows()
	for rows.Next() {
		var records OrderTable
		err = rows.Scan(&records.Id, &records.SN, &records.Mid, &records.MType, &records.Tokenid, &records.CreateTime,
			&records.OnchainTime, &records.Txhash, &records.Buyer, &records.Seller, &records.Price,
			&records.TokenType, &records.Status, &records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("OrderGetByOid error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func OrdersByTokenId(tokenid, mid, ismid, status, chain_type string) (recordList []OrderTable) {
	var rows *sql.Rows
	var err error
	var sql string

	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf("And chain_type = '%s'", chain_type)
	}

	if ismid == "0" { // 0 是不包含； 1 是包含
		sql = fmt.Sprintf("SELECT id,sn,mid,mtype,tokenid,create_time,onchain_time,txhash,"+
			"buyer,seller,price,token_type,status, free_gas, donation FROM t_order WHERE tokenid = '%s' AND mid != '%s' AND status = '%s' %s limit 10", tokenid, mid, status, chain_type_sql)
	} else {
		sql = fmt.Sprintf("SELECT id,sn,mid,mtype,tokenid,create_time,onchain_time,txhash,"+
			"buyer,seller,price,token_type,status, free_gas, donation FROM t_order WHERE tokenid = '%s' AND mid = '%s' AND status = '%s' %s limit 10", tokenid, mid, status, chain_type_sql)
	}
	rows, err = Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("OrdersByTokenId error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records OrderTable
		err = rows.Scan(&records.Id, &records.SN, &records.Mid, &records.MType, &records.Tokenid, &records.CreateTime,
			&records.OnchainTime, &records.Txhash, &records.Buyer, &records.Seller, &records.Price,
			&records.TokenType, &records.Status, &records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("OrdersByTokenId error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func OrdersGetMaxPrice(tokenid, mid, chain_type string) (price float32) {
	var rows *sql.Rows
	var err error
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf("and chain_type = '%s'", chain_type)
	}
	sql := fmt.Sprintf("SELECT SUM(cast (price as float8)) AS total FROM t_order WHERE tokenid = '%s' AND mid = %s "+
		" and status = 1 %s GROUP BY buyer ORDER BY total DESC LIMIT 1", tokenid, mid, chain_type_sql)
	rows, err = Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("OrdersGetMaxPrice error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&price)
		if err != nil {
			logger.Error("OrdersGetMaxPrice error", "err", err)
		}
	}
	return price
}

// OrdersByBuyer 获取购买者交易
func OrdersByBuyer(buyer, mtype, status, chain_type string) (recordList []OrderTable) {
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf("and chain_type = '%s'", chain_type)
	}

	sql := fmt.Sprintf("SELECT DISTINCT on (tokenid,chain_type) id,sn,mid,mtype,tokenid,create_time,onchain_time,txhash,"+
		"buyer,seller,price,token_type,status, chain_type, free_gas, donation FROM t_order WHERE buyer = '%s' AND mtype = '%s' AND"+
		" status = '%s' %s ", buyer, mtype, status, chain_type_sql)
	rows, err := Db.Raw(sql).Rows()
	for rows.Next() {
		var records OrderTable
		err = rows.Scan(&records.Id, &records.SN, &records.Mid, &records.MType, &records.Tokenid, &records.CreateTime,
			&records.OnchainTime, &records.Txhash, &records.Buyer, &records.Seller, &records.Price,
			&records.TokenType, &records.Status, &records.ChainType, &records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("OrdersByBuyer error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// OrdersCountByMid 获取指定mid的条数
func OrdersCountByMid(mid int) int {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM t_order WHERE mid = '%s' AND mtype = 1 AND status = 1", mid)
	rows, err := Db.Raw(sql).Rows()
	total := 0
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			logger.Error("OrdersCountByMid error", "err", err)
		}
	}
	return total
}
