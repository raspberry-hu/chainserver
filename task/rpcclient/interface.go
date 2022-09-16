package rpcclient

// EthRPCAPI 以太坊RPC接口
type EthRPCAPI interface {
	GetLatestBlock(isFullTx bool) (interface{}, error)
	GetBlockByNumber(height uint64, isFullTx bool) (interface{}, error)
	GetBlockByHash(hash string, isFullTx bool) (interface{}, error)
	GetTransactionByHash(hash string) (interface{}, error)
	GetTransactionReceipt(hash string) (interface{}, error)
	GetCode(contract string, height uint64) (interface{}, error)
}

// AssetsManagePushAPI 资产管理推送接口
type AssetsManagePushAPI interface {
	PushTransaction(hash string, status uint64) error
}
