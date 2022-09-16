package response

type NftView struct {
	NftId        int    `json:"nft_id"`        // 主键id
	NftName      string `json:"nft_name"`      // 名称
	NftDesc      string `json:"nft_desc"`      // 描述
	RightsRules  string `json:"rights_rules"`  // 权益规则
	TokenId      string `json:"token_id"`      // nft编号
	TokenUri     string `json:"token_uri"`     // nft地址
	Txhash       string `json:"txhash"`        // 铸造nft的hash
	Creater      int    `json:"creater"`       // nft创建者
	CreateNumber int    `json:"create_number"` // 创建区块高度
	CreateTime   string `json:"create_time"`   // 创建时间戳
	MediaUri     string `json:"media_uri"`     // 数字资产地址
	CreateTax    int    `json:"create_tax"`    // 铸造税
	Owner        int    `json:"owner"`         // nft拥有者
	NftType      int    `json:"nft_type"`      // nft分类
	Status       int    `json:"status"`        // nft铸造状态
	MarKetStatus int    `json:"market_status"` // nft挂单状态
	ChainType    string `json:"chain_type"`    // 链类型
	FreeGas      string `json:"free_gas"`      // 是否免收gas费 0 否 1 是
	Donation     int    `json:"donation"`      // 是否捐赠 0 否 1 是
	CollectionId int    `json:"collection_id"` // 所属的集合id
}

type NftAndMarketView struct {
	NftId            int    `json:"nft_id"`          // 主键id
	SN               string `json:"sn"`              // 流水号
	NftName          string `json:"nft_name"`        // 名称
	NftDesc          string `json:"nft_desc"`        // 描述
	TokenId          string `json:"token_id"`        // nft编号
	TokenUri         string `json:"token_uri"`       // nft地址
	Txhash           string `json:"txhash"`          // 铸造nft的hash
	Creater          int    `json:"creater"`         // nft创建者
	CreateNumber     int    `json:"create_number"`   // 创建区块高度
	CreateTime       string `json:"create_time"`     // 创建时间戳
	MediaUri         string `json:"media_uri"`       // 数字资产地址
	MediaIpfsUri     string `json:"media_ipfs_uri"`  // 图片ipfs地址
	ExploreUri       string `json:"explore_uri"`     // 音视频资产的展示图片
	CreateTax        int    `json:"create_tax"`      // 铸造税
	Owner            int    `json:"owner"`           // nft拥有者
	NftType          int    `json:"nft_type"`        // nft分类
	CurrencyName     string `json:"currency_name"`   // 币种的名称
	Status           int    `json:"status"`          // nft铸造状态
	MarKetStatus     int    `json:"market_status"`   // nft挂单状态
	CollectionId     int    `json:"collection_id"`   // 所属的集合id
	FavoritesCount   string `json:"favorites_count"` // 被收藏的个数
	MId              int    `json:"mid"`
	MarketCreateTime string `json:"market_create_time"` // 挂单时间
	StartingPrice    string `json:"starting_price"`     // 起拍价
	TokenType        string `json:"token_type"`         // 币种：BNB、USDT、CACA
	EndTime          string `json:"end_time"`           // 拍卖结束时间
	Bonus            int    `json:"bonus"`              // 拍卖参与者分红
	Attention        int    `json:"attention"`          // 是否关注
	ChainType        string `json:"chain_type"`         // 链类型
	OrderStatus      int    `json:"order_status"`       // 交易状态 对应market_list status 字段 0 已提交 1 已上链 2 失败 3 取消挂单 4 挂单成交
	FreeGas          string `json:"free_gas"`           // 是否免收gas费 0 否 1 是
	Donation         int    `json:"donation"`           // 是否捐赠 0 否 1 是
	OwnerUserName    string `json:"owner_user_name"`    // 拥有者用户名
	OwnerImageUrl    string `json:"owner_image_url"`    // 拥有者头像地址
	UserName         string `json:"user_name"`          // 创建者用户名
	ImageUrl         string `json:"image_url"`          // 创建者头像地址

}

type NftType struct {
	Id       int    `json:"id"`
	TypeName string `json:"type_name"` // 分类名称
}
