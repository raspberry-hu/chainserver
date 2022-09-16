package response

type MarketView struct {
	Creater       string `json:"creater"`        // 挂单者
	Tokenid       string `json:"tokenid"`        // 挂单nft编号
	MarketType    int    `json:"market_type"`    // 挂单类型：0 限价购买；1 拍卖
	StartingPrice string `json:"starting_price"` // 起拍价
	TokenType     string `json:"token_type"`     // 币种：BNB、USDT、CACA
	EndTime       string `json:"end_time"`       // 拍卖结束时间
	Buyer         string `json:"buyer"`          // 拍卖参与者分红
	Bonus         int    `json:"bonus"`          // 拍卖参与者分红
	Txhash        string `json:"txhash"`         // 交易hash
	CreateTime    string `json:"create_time"`    // 挂单时间
	Sorting       int    `json:"sorting"`        // 排序
	Status        int    `json:"status"`         // 挂单状态：0 开启 1 关闭
	FreeGas       string `json:"free_gas"`       // 是否免收gas费 0 否 1 是
	Donation      string `json:"donation"`       // 是否捐赠 0 否 1 是
	NftId         int    `json:"nft_id"`         // nft id

}

type TokenView struct {
	NftName        string `json:"nft_name"`        // 名称 ===> tokenId
	NftDesc        string `json:"nft_desc"`        // 描述 ===> 拥有者
	Creater        int    `json:"creater"`         // nft创建者
	CreaterName    string `json:"creater_name"`    // 创建者名字
	CreateTime     int64  `json:"create_time"`     // 铸造的时间
	ImageUrl       string `json:"image_url"`       // 创建者头像地址
	MarketType     int    `json:"market_type"`     // 挂单类型：2 限价购买；1 拍卖
	MediaUri       string `json:"media_uri"`       // 数字资产地址 ===> hash值
	Owner          int    `json:"owner"`           // nft拥有者
	OwnerName      string `json:"owner_name"`      // 拥有者名字
	OwnerImageUrl  string `json:"owner_image_url"` // 拥有者头像地址
	BuyTime        int    `json:"buy_time"`        // 购买的时间
	MetaDataUri    string `json:"meta_data_uri"`   // 原数据文件token url
	CollectionId   int    `json:"collection_id"`   // 集合id
	CollectionName string `json:"collection_name"` // 集合名字
	ExploreUri     string `json:"explore_uri"`     // 音视频资产用于展示的图片地址
	AntTokenId     int    `json:"ant_token_id"`    // 公链上的tokenId
	AntNftId       int    `json:"ant_nft_id"`      // 蚂蚁链中当前资产的Id
	AntCount       int    `json:"ant_count"`       // 蚂蚁链的铸造个数
	AntPrice       int    `json:"ant_price"`       // 资产价格
	AntNftOwner    string `json:"ant_nft_owner"`   // 蚂蚁链的nft铸造者
	AntTokenUrl    string `json:"ant_token_url"`   // 蚂蚁链的tokenUrl
	AntTxHash      string `json:"ant_tx_hash"`     // 蚂蚁链藏品对应公链的hash值
	TransferHash   string `json:"transfer_hash"`   // 蚂蚁链进行交易的hash值
}

type MarketAndNftView struct {
	Id             int    `json:"id"`
	SN             string `json:"sn"`              // 流水号
	Creater        string `json:"creater"`         // 挂单者
	UserName       string `json:"user_name"`       // 创建者用户名
	ImageUrl       string `json:"image_url"`       // 创建者头像地址
	OwnerUserName  string `json:"owner_user_name"` // 拥有者用户名
	OwnerImageUrl  string `json:"owner_image_url"` // 拥有者头像地址
	Tokenid        string `json:"tokenid"`         // 挂单nft编号
	MarketType     int    `json:"market_type"`     // 挂单类型：0 限价购买；1 拍卖
	StartingPrice  string `json:"starting_price"`  // 起拍价
	TokenType      string `json:"token_type"`      // 币种：BNB、USDT、CACA
	EndTime        string `json:"end_time"`        // 拍卖结束时间
	Buyer          string `json:"buyer"`           // 拍卖参与者分红
	Bonus          int    `json:"bonus"`           // 拍卖参与者分红
	Txhash         string `json:"txhash"`          // 交易hash
	CancelHash     string `json:"cancel_hash"`     // 取消挂单hash
	CreateTime     string `json:"create_time"`     // 挂单时间
	Sorting        int    `json:"sorting"`         // 排序
	Status         int    `json:"status"`          // 挂单状态：0 开启 1 关闭
	NftName        string `json:"nft_name"`        // 名称
	NftDesc        string `json:"nft_desc"`        // 描述
	TokenId        string `json:"token_id"`        // nft编号
	TokenUri       string `json:"token_uri"`       // nft地址
	NftTxhash      string `json:"nft_txhash"`      // 铸造nft的hash
	NftCreater     int    `json:"nft_creater"`     // nft创建者
	CreateNumber   int    `json:"create_number"`   // 创建区块高度
	Nft_CreateTime int64  `json:"nft_create_time"` // 创建时间戳
	MediaUri       string `json:"media_uri"`       // 数字资产地址
	CreateTax      int    `json:"create_tax"`      // 铸造税
	Owner          int    `json:"owner"`           // nft拥有者
	NftType        int    `json:"nft_type"`        // nft分类
	Approved       int    `json:"approved"`        // nft审核状态
	OfferCount     int    `json:"offer_count"`     // 竞拍出价次数
	Attention      int    `json:"attention"`       // 是否关注
	AttentionCount int    `json:"attention_count"` // 关注次数
	ChainType      string `json:"chain_type"`      // 链类型
	FreeGas        string `json:"free_gas"`        // 是否免收gas费 0 否 1 是
	Donation       string `json:"donation"`        // 是否捐赠 0 否 1 是
	NftId          int    `json:"nft_id"`          // nft id
}
