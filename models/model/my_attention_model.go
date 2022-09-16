package model

import (
	. "ChainServer/models/postgresql"
	"fmt"
	"github.com/sea-project/go-logger"
)

func MyAttentionInsert(wallet_addr, tokenid, chain_type string, nft_id int) error {
	sql := "INSERT INTO t_my_attention (wallet_addr,tokenid,chain_type,nft_id) VALUES('%s','%s','%s',%d)"
	sql = fmt.Sprintf(sql, wallet_addr, tokenid, chain_type, nft_id)
	Db.Exec(sql)
	return nil
}

func MyAttentionFindNew(wallet_addr string, nft_id int) string {
	var id string
	sql := "SELECT id FROM t_my_attention where wallet_addr = '%s' and nft_id = %d"
	sql = fmt.Sprintf(sql, wallet_addr, nft_id)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("MyAttentionFind error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			logger.Error("MyAttentionFind error", "err", err)
		}

	}
	return id
}

func MyAttentionDel(nft_id int) error {
	sql := "DELETE FROM t_my_attention where nft_id = %d"
	sql = fmt.Sprintf(sql, nft_id)
	Db.Exec(sql)
	return nil
}

func MyAttentionCount(wallet_addr string) int {
	var count int
	sql := "SELECT count(*) FROM t_my_attention where wallet_addr = '%s'"
	sql = fmt.Sprintf(sql, wallet_addr)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("MyAttentionFind error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			logger.Error("MyAttentionFind error", "err", err)
		}
	}
	return count
}
