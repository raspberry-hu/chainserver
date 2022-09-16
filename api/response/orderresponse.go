package response

type OrderView struct {
	Mid         int    `json:"mid"`          // 对应的挂单id
	Tokenid     string `json:"tokenid"`      // tokenid
	CreateTime  string `json:"create_time"`  // 下单时间
	OnchainTime string `json:"onchain_time"` // 上链时间
	Txhash      string `json:"txhash"`       // 交易hash
	Buyer       string `json:"buyer"`        // 买家
	Seller      string `json:"seller"`       // 卖家
	Price       string `json:"price"`        // 价格
	TokenType   string `json:"token_type"`   // 币种：BNB、USDT、CACA
	Status      int    `json:"status"`       // 订单状态：0 已提交 1 已上链 2 失败
	FreeGas     string `json:"free_gas"`     // 是否免收gas费 0 否 1 是
	Donation    string `json:"donation"`     // 是否捐赠 0 否 1 是
}

type MyBids struct {
	Mid          int    `json:"mid"`           // 对应的挂单id
	Tokenid      string `json:"tokenid"`       // tokenid
	CreateTime   string `json:"create_time"`   // 下单时间
	OnchainTime  string `json:"onchain_time"`  // 上链时间
	EndTime      string `json:"end_time"`      // 拍卖结束时间
	Txhash       string `json:"txhash"`        // 交易hash
	Buyer        string `json:"buyer"`         // 买家
	Seller       string `json:"seller"`        // 卖家
	Price        string `json:"price"`         // 价格
	TokenType    string `json:"token_type"`    // 币种：BNB、USDT、CACA
	Status       int    `json:"status"`        // 订单状态：0 已提交 1 已上链 2 失败
	NftName      string `json:"nft_name"`      // 名称
	NftDesc      string `json:"nft_desc"`      // 描述
	MediaUri     string `json:"media_uri"`     // 图片
	CreateNumber int    `json:"create_number"` // 区块高度
	Attention    int    `json:"attention"`     // 是否关注
	ChainType    string `json:"chain_type"`    // 链类型
	OrderStatus  int    `json:"order_status"`  // 对应 market_list status 0 已提交 1 已上链 2 失败 3 取消挂单 4 挂单成交',
	FreeGas      string `json:"free_gas"`      // 是否免收gas费 0 否 1 是
	Donation     string `json:"donation"`      // 是否捐赠 0 否 1 是
}
