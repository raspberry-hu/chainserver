package request

// sn,mid,tokenid,create_time,buyer,seller,price,token_type
type OrderNew struct {
	MarketId     int    `json:"market_id"`     // 对应的挂单id
	NftId        int    `json:"nft_id"`        // 对应的nft id
	MarketType   int    `json:"market_type"`   // 对应的挂单类型
	TokenId      string `json:"token_id"`      // tokenid
	CreateTime   int    `json:"create_time"`   // 下单时间
	Buyer        int    `json:"buyer"`         // 买家
	Seller       int    `json:"seller"`        // 卖家
	Price        string `json:"price"`         // 价格
	CurrencyName string `json:"currency_name"` // 币种：BNB、USDT、CACA
	ChainName    string `json:"chain_name"`    // 链类型
	Lazy         int    `json:"lazy"`          // 是否免收gas费 0 否 1 是
	Donation     int    `json:"donation"`      // 是否捐赠 0 否 1 是
	TxHash       string `json:"tx_hash"`       // 是否捐赠 0 否 1 是
}
