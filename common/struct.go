package common

type Request struct {
	ID      string   `json:"id"`
	Mthd    string   `json:"method"`
	Args    []string `json:"params"`
	Version string   `json:"jsonrpc"`
}

type Response struct {
	ID      string      `json:"id"`
	Version string      `json:"jsonrpc"`
	Code    int         `json:"code"`
	Error   *Error      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type ResponseNew struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result,omitempty"`
	Count  int         `json:"count"`
}

type ResponseCommon struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result,omitempty"`
}

type Error struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
	Data string `json:"data,omitempty"`
}

type ApiTransaction struct {
	BlockNumber       string `json:"blockNumber,omitempty"`
	TimeStamp         string `json:"timeStamp,omitempty"`
	Hash              string `json:"hash,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	BlockHash         string `json:"blockHash,omitempty"`
	TransactionIndex  string `json:"transactionIndex,omitempty"`
	From              string `json:"from,omitempty"`
	To                string `json:"to,omitempty"`
	Value             string `json:"value,omitempty"`
	Gas               string `json:"gas,omitempty"`
	GasPrice          string `json:"gasPrice,omitempty"`
	IsError           string `json:"isError,omitempty"`
	TxreceiptStatus   string `json:"txreceipt_status,omitempty"`
	Input             string `json:"input,omitempty"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	CumulativeGasUsed string `json:"cumulativeGasUsed,omitempty"`
	GasUsed           string `json:"gasUsed,omitempty"`
	Confirmations     string `json:"confirmations,omitempty"`
}
