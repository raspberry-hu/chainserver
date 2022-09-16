package request

// sn,creater,tokenid,market_type,starting_price,token_type,end_time,bonus
type MarketNew struct {
	NftId          int    `json:"nft_id"`         // nft 主键id
	Creater        int    `json:"creater"`        // 挂单者
	TokenId        string `json:"token_id"`       // 挂单nft编号
	MarketType     int    `json:"market_type"`    // 挂单类型：0 限价购买；1 拍卖
	StartingPrice  string `json:"starting_price"` // 起拍价
	EndTime        int    `json:"end_time"`       // 拍卖结束时间
	Buyer          int    `json:"buyer"`
	Reward         int    `json:"reward"`          // 拍卖参与者分红
	MetaDataUri    string `json:"meta_data_uri"`   // 原数据文件 ipfsurl
	CurrencyName   string `json:"currency_name"`   // 币种
	Lazy           int    `json:"lazy"`            // 是否免收gas费 0 否 1 是
	Donation       int    `json:"donation"`        // 是否捐赠 0 否 1 是
	DonationUserId int    `json:"donation_userid"` // 捐赠地址
	ChainName      string `json:"chain_name"`      // 链类型
	TxHash         string `json:"tx_hash"`         // 链类型
}
type TMarketList struct {
	Id            int
	NftId         int
	Creater       int
	TokenId       string
	MarketType    int
	StartingPrice string
	EndTime       int

	Buyer          int
	Reward         int
	TxHash         string
	CancelHash     string
	DealHash       string
	CreateTime     int
	Status         int
	ChainName      string
	CurrencyName   string
	Lazy           int
	Donation       int
	DonationUserId int
}
