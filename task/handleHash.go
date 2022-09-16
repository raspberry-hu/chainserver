package task

import (
	"ChainServer/common"
	"ChainServer/task/rpcclient"
	"ChainServer/task/rpcclient/jsonprc"
	"github.com/sea-project/go-logger"
	"strconv"
	"strings"
)

type ChainInfo struct {
	TokenId      string `:"token_id"`      // nft编号
	CreateNumber string `:"create_number"` // 上链区块高度
	CreateTime   string `:"create_time"`   // 上链时间
	Owner        string `:"owner"`         // nft拥有者
}

// handleNftHash return token_id,create_number,create_time,owner
func handleNftHash(client *jsonrpc.Http, hash string) (ChainInfo, error) {
	// 根据hash获取区块高度
	tx, err := rpcclient.GetOnChainTxInfo(client, hash)
	if err != nil {
		return ChainInfo{}, err
	}
	logger.Info("根据hash获取区块高度 hash === ", hash, "tx =====", tx)
	// 根据区块高度获取区块时间
	block, err := rpcclient.GetOnChainBlockInfo(client, tx.BlockNumber)
	logger.Info("根据区块高度获取区块时间 hash === ", hash, "block =====", block)

	if err != nil {
		return ChainInfo{}, err
	}
	// 根据hash获取交易回执，解析交易内nft授权信息
	receipt, err := rpcclient.GetOnChainReceiptInfo(client, hash)
	logger.Info("根据hash获取交易回执，解析交易内nft授权信息 hash === ", hash, "receipt =====", receipt)

	if err != nil {
		return ChainInfo{}, err
	}
	// 组装信息
	var token_id string
	if len(receipt.Value) > 2 {
		receiptValue, err := common.HexToInt64(receipt.Value)
		if err != nil {
			return ChainInfo{}, err
		}
		token_id = strconv.FormatInt(receiptValue, 10)
	} else {
		token_id = ""
	}

	blockNumber, err := common.HexToInt64(tx.BlockNumber)
	if err != nil {
		return ChainInfo{}, err
	}
	create_number := strconv.FormatInt(blockNumber, 10)
	blockTimestamp, err := common.HexToInt64(block.Timestamp)
	if err != nil {
		return ChainInfo{}, err
	}
	//create_time := common.IntToDate(blockTimestamp)
	create_time := strconv.FormatInt(blockTimestamp, 10)
	receiptTo := rpcclient.BytesToAddress(rpcclient.FromHex(receipt.To))
	if err != nil {
		return ChainInfo{}, err
	}
	owner := strings.ToLower(receiptTo.Hex())

	return ChainInfo{TokenId: token_id, CreateNumber: create_number, CreateTime: create_time, Owner: owner}, nil
}
