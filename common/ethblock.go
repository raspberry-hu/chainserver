package common

// 全交易区块
type Block struct {
	Number     string         `json:"number"`
	Size       string         `json:"size"`
	Timestamp  string         `json:"timestamp"`
	GasLimit   string         `json:"gasLimit"`
	GasUsed    string         `json:"gasUsed"`
	Hash       string         `json:"hash"`
	ParentHash string         `json:"parentHash"`
	Coinbase   string         `json:"miner"`
	TxDatas    []*Transaction `json:"transactions"`
}

// 带哈希列表的区块
type BlockTxHash struct {
	Number     string   `json:"number"`
	Size       string   `json:"size"`
	Timestamp  string   `json:"timestamp"`
	GasLimit   string   `json:"gasLimit"`
	GasUsed    string   `json:"gasUsed"`
	Hash       string   `json:"hash"`
	ParentHash string   `json:"parentHash"`
	Coinbase   string   `json:"miner"`
	TxHexs     []string `json:"transactions"`
}

type Header struct {
	Number     string `json:"number"`
	Size       string `json:"size"`
	Timestamp  string `json:"timestamp"`
	GasLimit   string `json:"gasLimit"`
	GasUsed    string `json:"gasUsed"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Coinbase   string `json:"miner"`
}

type Transaction struct {
	BlockHash   string `json:"blockHash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	Recipient   string `json:"to"`
	GasLimit    string `json:"gas"`
	GasPrice    string `json:"gasPrice"`
	Hash        string `json:"hash"`
	Payload     string `json:"input"`
	Nonce       string `json:"nonce"`
	R           string `json:"r"`
	S           string `json:"s"`
	V           string `json:"v"`
	Index       string `json:"transactionIndex"`
	Value       string `json:"value"`
}

type TxReceipt struct {
	Hash         string         `json:"transactionHash"`
	Index        string         `json:"transactionIndex"`
	BlockNumber  string         `json:"blockNumber"`
	BlockHash    string         `json:"blockHash"`
	GasUsedTotal string         `json:"cumulativeGasUsed"`
	GasUsed      string         `json:"gasUsed"`
	Contract     string         `json:"contractAddress"`
	LogsBloom    string         `json:"logsBloom"`
	Status       string         `json:"status"`
	From         string         `json:"from"`
	To           string         `json:"to"`
	Logs         []*ReceiptLogs `json:"logs"`
}

type ReceiptLogs struct {
	Address     string   `json:"address"`
	BlockNumber string   `json:"blockNumber"`
	BlockHash   string   `json:"blockHash"`
	Index       string   `json:"transactionIndex"`
	Hash        string   `json:"transactionHash"`
	Data        string   `json:"data"`
	LogIndex    string   `json:"logIndex"`
	Removed     bool     `json:"removed"`
	Topics      []string `json:"topics"`
}
