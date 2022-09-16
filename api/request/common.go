package request

type UpdateHash struct {
	TxHash       string `json:"tx_hash"`
	NftId        int    `json:"nft_id"`
	MarketId     int    `json:"market_id"`
	OrderId      int    `json:"order_id"`
	TransferHash string `json:"transfer_hash"`
}

type UpdateCancelHash struct {
	CancelHash string `json:"cancel_hash"`
	DealHash   string `json:"deal_hash"`
	Lazy       int    `json:"lazy"`
	MarketId   int    `json:"market_id"`
	NftId      int    `json:"nft_id"`
	TokenId    string `json:"token_id"`
}
