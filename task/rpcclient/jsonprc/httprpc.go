package jsonrpc

import (
	"ChainServer/common"
	"ChainServer/task/rpcclient/http"
	"encoding/json"
	"errors"
	"github.com/sea-project/go-logger"
	"sync"
)

type Http struct {
	rpc *http.Client
}

var (
	httpRPC  *Http
	httpOnce sync.Once
)

var args []interface{}

// NewETHTP
func NewETHTP(host string) *Http {
	/*
		httpOnce.Do(func() {
			httpRPC = &Http{rpc: http.NewClient(host)}
		})*/
	httpRPC = &Http{rpc: http.NewClient(host)}
	return httpRPC
}

// NewETHTP2
func NewETHTP2(host string) *Http {
	return &Http{rpc: http.NewClient(host)}
}

// GetLatestBlock 获取最新区块信息,isFullTx 全交易
func (eth *Http) GetLatestBlock(isFullTx bool) (interface{}, error) {
	args = []interface{}{"latest", isFullTx}
	params := NewHttpParams("eth_getBlockByNumber", args)
	resBody, err := eth.rpc.HttpRequest(params)
	if err != nil {
		return nil, err
	}
	return eth.ParseJsonRPCResponse(resBody)
}

// GetBlockByNumber 根据区块高度获取区块信息,isFullTx 全交易
func (eth *Http) GetBlockByNumber(height string, isFullTx bool) (interface{}, error) {
	args = []interface{}{height, isFullTx}
	params := NewHttpParams("eth_getBlockByNumber", args)
	resBody, err := eth.rpc.HttpRequest(params)
	if err != nil {
		return nil, err
	}
	return eth.ParseJsonRPCResponse(resBody)
}

// GetBlockByHash 根据哈希获取区块信息，isFullTx 全交易
func (eth *Http) GetBlockByHash(hash string, isFullTx bool) (interface{}, error) {
	args = []interface{}{hash, isFullTx}
	params := NewHttpParams("eth_getBlockByHash", args)
	resBody, err := eth.rpc.HttpRequest(params)
	if err != nil {
		return nil, err
	}
	return eth.ParseJsonRPCResponse(resBody)
}

// GetTransactionByHash 获取交易信息
func (eth *Http) GetTransactionByHash(hash string) (interface{}, error) {
	args = []interface{}{hash}
	params := NewHttpParams("eth_getTransactionByHash", args)
	resBody, err := eth.rpc.HttpRequest(params)
	if err != nil {
		return nil, err
	}
	return eth.ParseJsonRPCResponse(resBody)
}

// GetTransactionReceipt 获取交易票据
func (eth *Http) GetTransactionReceipt(hash string) (interface{}, error) {
	args = []interface{}{hash}
	params := NewHttpParams("eth_getTransactionReceipt", args)
	resBody, err := eth.rpc.HttpRequest(params)
	if err != nil {
		return nil, err
	}
	return eth.ParseJsonRPCResponse(resBody)
}

// GetCode 获取链码
func (eth *Http) GetCode(contract string, height uint64) (interface{}, error) {
	tag := "latest"
	if height > 0 {
		tag = common.IntToHex(height)
	}
	args = []interface{}{contract, tag}
	params := NewHttpParams("eth_getCode", args)
	resBody, err := eth.rpc.HttpRequest(params)
	if err != nil {
		return nil, err
	}
	return eth.ParseJsonRPCResponse(resBody)
}

// ParseJsonRPCResponse jsonrpc格式返回结果解析
func (eth *Http) ParseJsonRPCResponse(resBody []byte) (interface{}, error) {
	response := new(common.Response)
	err := json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		logger.Error("---++++", "resBody", string(resBody))
		return nil, errors.New(string(resBody))
	}
	return response.Result, nil
}

type Request struct {
	ID      string        `json:"id"`
	Mthd    string        `json:"method"`
	Args    []interface{} `json:"params"`
	Version string        `json:"jsonrpc"`
}

func NewHttpParams(method string, args []interface{}) string {
	id := common.GetRandString(16)
	request := &Request{
		ID:      id,
		Mthd:    method,
		Args:    args,
		Version: "2.0",
	}
	rb, _ := json.Marshal(request)
	return string(rb)
}
