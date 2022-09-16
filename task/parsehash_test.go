package task

import (
	"ChainServer/common"
	"ChainServer/task/rpcclient"
	"ChainServer/task/rpcclient/jsonprc"
	"encoding/json"
	"testing"
)

func Test_gettime(t *testing.T) {
	client := jsonrpc.NewETHTP("config.Rpc")
	res, err := client.GetBlockByNumber("0x123", true)
	if err != nil {
		t.Log(err)
	}
	block := new(common.Block)
	resB, err := json.Marshal(res)
	if err != nil {
		t.Log(err)
	}
	if err := json.Unmarshal(resB, &block); err != nil {
		t.Log(err)
	}
	t.Log(block.Timestamp)
	num, err := common.HexToInt64(block.Timestamp)
	if err != nil {
		t.Log(err)
	}
	t.Log(num)
}

func Test_gettx(t *testing.T) {
	client := jsonrpc.NewETHTP("config.Rpc")
	res, err := client.GetTransactionByHash("0xd89759fe997116af9ccf10b814d02bfe7579f081de303e50e623447cf1435a62")
	if err != nil {
		t.Log(err)
	}
	transaction := new(common.Transaction)
	resB, err := json.Marshal(res)
	if err != nil {
		t.Log(err)
	}
	if err := json.Unmarshal(resB, &transaction); err != nil {
		t.Log(err)
	}
	t.Log(transaction.BlockNumber)
}

func Test_getrep(t *testing.T) {
	client := jsonrpc.NewETHTP("config.Rpc")
	res, err := client.GetTransactionReceipt("0xd89759fe997116af9ccf10b814d02bfe7579f081de303e50e623447cf1435a62")
	if err != nil {
		t.Log(err)
	}
	txReceipt := new(common.TxReceipt)
	resB, err := json.Marshal(res)
	if err != nil {
		t.Log(err)
	}
	if err := json.Unmarshal(resB, &txReceipt); err != nil {
		t.Log(err)
	}
	t.Log(txReceipt.Status)

	/*receiptLogs := new(common.ReceiptLogs)
	resR, err := json.Marshal(txReceipt.Logs)
	if err != nil {
		t.Log(err)
	}
	if err := json.Unmarshal(resR, &receiptLogs); err != nil {
		t.Log(err)
	}
	t.Log(receiptLogs.BlockNumber)*/
	for _, value := range txReceipt.Logs {
		if value.Address == "0x1d30eb002db028c8965f7b73329c9c296d7394c1" {
			t.Log(value.Topics)
		}
	}
}

func Test_GetOnChainTxInfo(t *testing.T) {
	client := jsonrpc.NewETHTP("config.Rpc")
	info, err := rpcclient.GetOnChainTxInfo(client, "0xd89759fe997116af9ccf10b814d02bfe7579f081de303e50e623447cf1435a62")
	if err != nil {
		t.Log(err)
	}
	t.Log(info.From, info.To, info.BlockNumber)
}

func Test_GetOnChainBlockInfo(t *testing.T) {
	client := jsonrpc.NewETHTP("config.Rpc")
	info, err := rpcclient.GetOnChainBlockInfo(client, "0x87a764")
	if err != nil {
		t.Log(err)
	}
	t.Log(info.Timestamp)
}

func Test_GetOnChainReceiptInfo(t *testing.T) {
	client := jsonrpc.NewETHTP("config.Rpc")
	info, err := rpcclient.GetOnChainReceiptInfo(client, "0xd89759fe997116af9ccf10b814d02bfe7579f081de303e50e623447cf1435a62")
	if err != nil {
		t.Log(err)
	}
	t.Log(info.ContractAddress, info.FuncName, info.From, info.To, info.Value)
	t.Log(common.HexToInt64(info.Value))
}

func Test_hex(t *testing.T) {
	a := "0x00000000000000000000000064cf0f4e7fae21f27ad08368256cfe85fb07bba6"
	b := rpcclient.BytesToAddress(rpcclient.FromHex(a))
	t.Log(b.Hex())
}
