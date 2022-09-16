package model

import (
	. "ChainServer/models/postgresql"
	"fmt"

	"github.com/sea-project/go-logger"
)

// &records.Id, &records.NftName, &records.NftDesc, &records.TokenId, &records.TokenUri,
// &records.Txhash, &records.Creater, &records.CreateNumber, &records.CreateTime, &records.MediaUri,
// &records.CreateTax, &records.Owner, &records.NftType, &records.Status
// `id`,`nft_name` ,`nft_desc`,`token_id`,`token_uri` ,`txhash` ,`creater`,`create_number`,`create_time`,
// `media_uri`,`create_tax`,`owner`,`nft_type`,`status`
type NftTable struct {
	Id           int    `:"id"`
	SN           string `:"sn"`                // 流水号
	NftName      string `:"nft_name"`          // 名称
	NftDesc      string `:"nft_desc"`          // 描述
	RightsRules  string `:"rights_rules"`      // 权益规则
	TokenId      string `:"token_id"`          // nft编号
	TokenUri     string `:"token_uri"`         // nft地址
	Txhash       string `:"txhash"`            // 铸造nft的hash
	Creater      int    `:"creater"`           // nft创建者
	CreateNumber int    `:"create_number"`     // 创建区块高度
	CreateTime   int64  `:"create_time"`       // 创建时间戳
	MediaUri     string `:"media_uri"`         // 数字资产地址
	CreateTax    int    `:"create_tax"`        // 铸造税
	Owner        int    `:"owner"`             // nft拥有者
	NftType      int    `:"nft_type"`          // nft分类
	Status       int    `:"status"`            // nft铸造状态
	MarKetStatus int    `json:"market_status"` // nft挂单状态
	Approved     int    `json:"approved"`      // nft审核状态
	ChainType    string `json:"chain_type"`    // 链类型
	TokenType    string `json:"token_type"`    // 币种类型
	FreeGas      string `json:"free_gas"`      // 是否免收gas费 0 否 1 是
	Donation     int    `json:"donation"`      // 是否捐赠 0 否 1 是
}

// GetNftHashList 获取有hash、没有更新链信息的记录
func GetNftHashList(chain_type string) (recordList []NftTable) {
	sql := "SELECT id,sn,nft_name ,nft_desc,rights_rules,token_id,token_uri ," +
		"txhash ,creater,create_number,create_time,media_uri,create_tax,owner," +
		"nft_type,status, market_status, approved FROM nft WHERE status = 0 and txhash !='' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, chain_type)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("GetNftHashList error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records NftTable
		err = rows.Scan(&records.Id, &records.SN, &records.NftName, &records.NftDesc, &records.RightsRules, &records.TokenId, &records.TokenUri,
			&records.Txhash, &records.Creater, &records.CreateNumber, &records.CreateTime, &records.MediaUri,
			&records.CreateTax, &records.Owner, &records.NftType, &records.Status, &records.MarKetStatus, &records.Approved)
		if err != nil {
			logger.Error("GetNftHashList error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// InsertRecord
func NFTInsertRecord(sn, nft_name, nft_desc, rights_rules, media_uri, create_tax, creater, owner, chain_type, token_type, free_gas, donation string, status int, ipfs_url, raw_info_url string, create_time int) error {
	sql := "INSERT INTO nft(sn,nft_name,nft_desc, rights_rules,media_uri," +
		"create_tax,creater,owner,chain_type,token_type,free_gas, donation, status, ipfs_uri, token_uri, create_time) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s', '%s', '%s', '%s', %d, '%s', '%s', %d)"
	sql = fmt.Sprintf(sql, sn, nft_name, nft_desc, rights_rules, media_uri, create_tax, creater, owner, chain_type,
		token_type, free_gas, donation, status, ipfs_url, raw_info_url, create_time)
	Db.Exec(sql)
	return nil
}

func NftGetByTokenId(tokenid, chain_type string) (recordList []NftTable) {
	var chain_type_sql string
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf("and chain_type = '%s'", chain_type)
	}
	sql := "SELECT id,sn,nft_name ,nft_desc,rights_rules,token_id,token_uri ," +
		"txhash ,creater,create_number,create_time,media_uri,create_tax,owner," +
		"nft_type,status, market_status, approved, chain_type, free_gas, donation FROM nft WHERE token_id = '%s' %s"
	sql = fmt.Sprintf(sql, tokenid, chain_type_sql)
	rows, err := Db.Raw(sql).Rows()
	defer rows.Close()
	for rows.Next() {
		var records NftTable
		err = rows.Scan(&records.Id, &records.SN, &records.NftName, &records.NftDesc, &records.RightsRules,
			&records.TokenId, &records.TokenUri, &records.Txhash, &records.Creater, &records.CreateNumber,
			&records.CreateTime, &records.MediaUri, &records.CreateTax, &records.Owner, &records.NftType, &records.Status,
			&records.MarKetStatus, &records.Approved, &records.ChainType, &records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("NftGetByTokenId error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func NftGetBySn(sn, chain_type string) (recordList []NftTable) {
	sql := "SELECT id,sn,nft_name ,nft_desc,rights_rules,token_id,token_uri ," +
		"txhash ,creater,create_number,create_time,media_uri,create_tax,owner," +
		"nft_type,status, market_status, approved, free_gas, donation FROM nft WHERE sn = '%s' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, sn, chain_type)
	rows, err := Db.Raw(sql).Rows()
	for rows.Next() {
		var records NftTable
		err = rows.Scan(&records.Id, &records.SN, &records.NftName, &records.NftDesc, &records.RightsRules,
			&records.TokenId, &records.TokenUri, &records.Txhash, &records.Creater, &records.CreateNumber,
			&records.CreateTime, &records.MediaUri, &records.CreateTax, &records.Owner, &records.NftType,
			&records.Status, &records.MarKetStatus, &records.Approved, &records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("NftGetBySn error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

func NftGetById(nft_id int) (recordList []NftTable) {
	sql := "SELECT id,sn,nft_name ,nft_desc,rights_rules,token_id,token_uri ," +
		"txhash ,creater,create_number,create_time,media_uri,create_tax,owner," +
		"nft_type,status, market_status, approved, free_gas, donation FROM nft WHERE id = %d"
	sql = fmt.Sprintf(sql, nft_id)
	rows, err := Db.Raw(sql).Rows()
	for rows.Next() {
		var records NftTable
		err = rows.Scan(&records.Id, &records.SN, &records.NftName, &records.NftDesc, &records.RightsRules,
			&records.TokenId, &records.TokenUri, &records.Txhash, &records.Creater, &records.CreateNumber,
			&records.CreateTime, &records.MediaUri, &records.CreateTax, &records.Owner, &records.NftType,
			&records.Status, &records.MarKetStatus, &records.Approved, &records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("NftGetBySn error", "err", err)
		}
		recordList = append(recordList, records)
	}
	return recordList
}

// NftUpdateMarketStatusByTokenid 根据tokenid更新状态
func NftUpdateMarketStatusByTokenid(tokenid, status, chain_type string) error {
	sql := "UPDATE nft SET market_status = '%s' WHERE token_id = '%s' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, status, tokenid, chain_type)
	Db.Exec(sql)
	return nil
}

// NftHashUpdate 更新铸造记录hash
func NftHashTokenIdUpdate(id int, token_id, create_number string) error {
	sql := "UPDATE nft SET token_id = %s, create_number = '%s' WHERE id = %d"
	sql = fmt.Sprintf(sql, token_id, create_number, id)
	Db.Exec(sql)
	return nil
}

// NftHashUpdate 更新铸造记录hash
func NftHashUpdate(sn, hash, chain_type string) error {
	sql := "UPDATE nft SET txhash = '%s' WHERE sn = '%s' and status = 0 and chain_type = '%s'"
	sql = fmt.Sprintf(sql, hash, sn, chain_type)
	Db.Exec(sql)
	return nil
}

// NftOwnerUpdate 更新nft 拥有者
func NftOwnerUpdate(owner, free_gas, donation string, nft_id int) error {
	sql := "UPDATE nft SET owner = '%s', free_gas = '%s', donation = '%s' WHERE id = %d"
	sql = fmt.Sprintf(sql, owner, free_gas, donation, nft_id)
	Db.Exec(sql)
	return nil
}

// NftStatusUpdate 更新nft 状态
func NftStatusUpdate(hash, chain_type string) error {
	sql := "UPDATE nft SET status = 2 WHERE txhash = '%s' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, hash, chain_type)
	Db.Exec(sql)
	return nil
}

// NftMarketStatusUpdate 更新nft market状态
func NftMarketStatusUpdate(tokenid string, market_status int, chain_type string) error {
	sql := "UPDATE nft SET market_status = '%d' WHERE token_id = '%s' and chain_type = '%s'"
	sql = fmt.Sprintf(sql, market_status, tokenid, chain_type)
	Db.Exec(sql)
	return nil
}

func NftMarketStatusUpdateById(market_status, nft_id int) error {
	sql := "UPDATE nft SET market_status = '%d' WHERE id = %d"
	sql = fmt.Sprintf(sql, market_status, nft_id)
	Db.Exec(sql)
	return nil
}

// NftChainUpdate 更新nft 上链信息: token_id,creater,create_number,create_time,owner
func NftChainUpdate(hash, token_id, creater, create_number, create_time, owner string) error {
	sql := "UPDATE nft SET token_id = '%s',creater = '%s',create_number = '%s'," +
		"create_time = '%s',owner = '%s',status = 1 WHERE txhash = '%s'"
	sql = fmt.Sprintf(sql, token_id, creater, create_number, create_time, owner, hash)
	Db.Exec(sql)
	return nil
}

func NftMySalability(address, market_status, chain_type string) (recordList []NftTable) {
	var chain_type_sql string
	sql := "SELECT id,sn,nft_name ,nft_desc,rights_rules,token_id,token_uri ," +
		"txhash ,creater,create_number,create_time,media_uri,create_tax,owner," +
		"nft_type,status,market_status, approved, chain_type, token_type, free_gas, donation FROM nft " +
		"WHERE owner = '%s' and status = 1  and market_status = '%s' "
	if chain_type != "" {
		chain_type_sql = fmt.Sprintf("and chain_type = '%s'", chain_type)
	}
	sql += chain_type_sql
	sql = fmt.Sprintf(sql, address, market_status)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("NftMySalability error", "err", err)
	}
	defer rows.Close()

	for rows.Next() {
		var records NftTable
		err = rows.Scan(&records.Id, &records.SN, &records.NftName, &records.NftDesc, &records.RightsRules, &records.TokenId, &records.TokenUri,
			&records.Txhash, &records.Creater, &records.CreateNumber, &records.CreateTime, &records.MediaUri,
			&records.CreateTax, &records.Owner, &records.NftType, &records.Status, &records.MarKetStatus, &records.Approved, &records.ChainType, &records.TokenType,
			&records.FreeGas, &records.Donation)
		if err != nil {
			logger.Error("NftMySalability error", "err", err)
		}
		//fmt.Println(records.CreateTime)
		recordList = append(recordList, records)
	}
	return recordList
}

// InsertNftChainInfo
func InsertNftChainInfo(hash, token_id, creater, create_number, create_time, owner, status string) error {
	sql := "UPDATE nft SET token_id = ?,creater = ?,create_number = ?,create_time = ?,owner = ?,status = ? WHERE txhash = ? and status = 0"
	Db.Exec(sql)
	return nil
}

// ------------------------------------------------------------------------------------------------------------

type Nft struct {
	ID           int    `json:"id"`
	Sn           string `json:"sn"`
	NftName      string `json:"nft_name"`
	NftDesc      string `json:"nft_desc"`
	RightsRules  string `json:"rights_rules"`
	TokenId      string `json:"token_id"`
	TokenUri     string `json:"token_uri"`
	TxHash       string `json:"tx_hash"`
	Creater      string `json:"creater"`
	CreateNumber int    `json:"create_number"`
	CreateTime   int    `json:"create_time"`
	MediaUri     string `json:"media_uri"`
	CreateTax    int    `json:"create_tax"`
	Owner        string `json:"owner"`
	NftType      int    `json:"nft_type"`
	Status       int    `json:"status"`
	MarketStatus int    `json:"market_status"`
	Approved     int    `json:"approved"`
	ChainType    string `json:"chain_type"`
	TokenType    string `json:"token_type"`
	FreeGas      string `json:"free_gas"`
	Donation     string `json:"donation"`
	IpfsUri      string `json:"ipfs_uri"`
}

func NftCommonFind(page, offset int, data interface{}) (nft_arry []Nft) {
	Db.Limit(offset).Offset(offset * (page - 1)).Where(data).Find(&nft_arry)
	return nft_arry
}
