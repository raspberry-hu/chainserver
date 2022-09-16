package rpcclient

import (
	"ChainServer/common"
	"ChainServer/config"
	"ChainServer/task/rpcclient/jsonprc"
	"encoding/json"
	"errors"
	"github.com/sea-project/go-logger"
	"strings"
)

// `token_id`,`creater`,`create_number`,`create_time`,`owner`
type NFTInfo struct {
	TokenId      string `json:"token_id"`
	Creater      string `json:"creater"`
	CreateNumber string `json:"create_number"`
	CreateTime   string `json:"create_time"`
	Owner        string `json:"owner"`
}

type TxInfo struct {
	From        string `json:"from"`
	To          string `json:"to"`
	BlockNumber string `json:"blockNumber"`
}

type BlockInfo struct {
	Timestamp string `json:"timestamp"`
}

type TransferReceiptInfo struct {
	ContractAddress string `json:"address"`
	FuncName        string `json:"func"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
}

// GetOnChainTxInfo 获取交易主要信息
func GetOnChainTxInfo(client *jsonrpc.Http, hash string) (TxInfo, error) {
	res, err := client.GetTransactionByHash(hash)
	if err != nil {
		logger.Error("GetOnChainTxInfo", "step", "GetTransactionByHash", "err", err)
		return TxInfo{}, err
	}

	transaction := new(common.Transaction)
	//logger.Info("transaction ==== ", transaction)
	resB, err := json.Marshal(res)
	if err != nil {
		logger.Error("GetOnChainTxInfo", "step", "Marshal res", "err", err)
		return TxInfo{}, err
	}
	if err := json.Unmarshal(resB, &transaction); err != nil {
		logger.Error("GetOnChainTxInfo", "step", "Unmarshal transaction", "err", err)
		return TxInfo{}, err
	}
	if transaction == nil {
		logger.Error("txhash not found", hash)
		return TxInfo{}, errors.New("txhash not found")
	}
	from := strings.ToLower(transaction.From)
	to := strings.ToLower(transaction.Recipient)
	blockNumber := transaction.BlockNumber
	return TxInfo{From: from, To: to, BlockNumber: blockNumber}, nil
}

// GetOnChainBlockInfo 获取区块主要信息
func GetOnChainBlockInfo(client *jsonrpc.Http, blockNumber string) (BlockInfo, error) {
	res, err := client.GetBlockByNumber(blockNumber, false)
	if err != nil {
		logger.Error("GetOnChainBlockInfo", "step", "GetBlockByNumber", "err", err)
		return BlockInfo{}, err
	}
	block := new(common.BlockTxHash)
	resB, err := json.Marshal(res)
	if err != nil {
		logger.Error("GetOnChainBlockInfo", "step", "Marshal res", "err", err)
		return BlockInfo{}, err
	}
	if err := json.Unmarshal(resB, &block); err != nil {
		logger.Error("GetOnChainBlockInfo", "step", "Unmarshal block", "err", err)
		return BlockInfo{}, err
	}

	return BlockInfo{Timestamp: block.Timestamp}, nil
}

// GetOnChainReceiptInfo 获取交易回执主要信息
func GetOnChainReceiptInfo(client *jsonrpc.Http, hash string) (TransferReceiptInfo, error) {
	res, err := client.GetTransactionReceipt(hash)
	if err != nil {
		logger.Error("GetOnChainBlockInfo", "step", "GetBlockByNumber", "err", err)
		return TransferReceiptInfo{}, err
	}
	txReceipt := new(common.TxReceipt)
	resB, err := json.Marshal(res)
	logger.Info("resB     ======    ", string(resB))
	if err != nil {
		logger.Error("GetOnChainBlockInfo", "step", "Marshal res", "err", err)
		return TransferReceiptInfo{}, err
	}
	if err := json.Unmarshal(resB, &txReceipt); err != nil {
		logger.Error("GetOnChainBlockInfo", "step", "Unmarshal block", "err", err)
		return TransferReceiptInfo{}, err
	}
	if txReceipt.Status != "0x1" {
		logger.Error("GetOnChainBlockInfo", "step", "txReceipt.Status", "hash", hash, "err", errors.New("contract vm running false"))
		return TransferReceiptInfo{}, errors.New("contract vm running false")
	}
	tri := TransferReceiptInfo{}
	// nft contract addr: 0x1d30eb002db028c8965f7b73329c9c296d7394c1
	// Transfer方法：0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
	for _, value := range txReceipt.Logs { // nftAddr
		logger.Info("jjjjjjjjjjjjjjjjjjjjjjjj", value.Address)
		logger.Info("vvvvvvvvvvvvvvvvvvvvvvvv", config.TopicList[0])

		if _, ok := config.NftAddrMap[value.Address]; ok && value.Topics[0] == config.TopicList[0] {
			logger.Info("mmmmmmmmmmmmmmmmmmmm")
			tri.ContractAddress = value.Address
			tri.FuncName = value.Topics[0]
			tri.From = strings.ToLower(value.Topics[1])
			tri.To = strings.ToLower(value.Topics[2])
			tri.Value = value.Topics[3]
		}
	}
	return tri, nil
}
