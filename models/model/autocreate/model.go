package autocreate

type Banner struct {
	Id         int    `:"id"`
	BannerName string `:"banner_name"` // banner名称
	BannerImg  string `:"banner_img"`  // banner图片地址
	BannerUrl  string `:"banner_url"`  // banner链接地址
}

type Market_list struct {
	Id            int    `:"id"`
	Creater       string `:"creater"`        // 挂单者
	Tokenid       string `:"tokenid"`        // 挂单nft编号
	MarketType    int    `:"market_type"`    // 挂单类型：0 限价购买；1 拍卖
	StartingPrice string `:"starting_price"` // 起拍价
	TokenType     string `:"token_type"`     // 币种：BNB、USDT、CACA
	EndTime       int    `:"end_time"`       // 拍卖结束时间
	Bonus         int    `:"bonus"`          // 拍卖参与者分红
	Txhash        string `:"txhash"`         // 交易hash
	CreateTime    int    `:"create_time"`    // 挂单时间
	Sorting       int    `:"sorting"`        // 排序
	Status        int    `:"status"`         // 状态：0 已提交 1 已上链 2 失败
}

type Nft struct {
	Id           int    `:"id"`
	NftName      string `:"nft_name"`      // 名称
	NftDesc      string `:"nft_desc"`      // 描述
	TokenId      string `:"token_id"`      // nft编号
	TokenUri     string `:"token_uri"`     // nft地址
	Txhash       string `:"txhash"`        // 铸造nft的hash
	Creater      string `:"creater"`       // nft创建者
	CreateNumber int    `:"create_number"` // 创建区块高度
	CreateTime   int    `:"create_time"`   // 创建时间戳
	MediaUri     string `:"media_uri"`     // 数字资产地址
	CreateTax    int    `:"create_tax"`    // 铸造税
	Owner        string `:"owner"`         // nft拥有者
	NftType      int    `:"nft_type"`      // nft分类
	Status       int    `:"status"`        // 状态：0 已提交 1 已上链 2 失败
}

type Nft_type struct {
	Id       int    `:"id"`
	TypeName string `:"type_name"` // 分类名称
}

type Order struct {
	Id          int    `:"id"`
	Mid         int    `:"mid"`          // 对应的挂单id
	Tokenid     string `:"tokenid"`      // tokenid
	CreateTime  int    `:"create_time"`  // 下单时间
	OnchainTime int    `:"onchain_time"` // 上链时间
	Txhash      string `:"txhash"`       // 交易hash
	Buyer       string `:"buyer"`        // 买家
	Seller      string `:"seller"`       // 卖家
	Price       string `:"price"`        // 价格
	TokenType   string `:"token_type"`   // 币种：BNB、USDT、CACA
	Status      int    `:"status"`       // 订单状态：0 已提交 1 已上链 2 失败
}

type Sys_info struct {
	Id    int    `:"id"`
	Key   string `:"key"`   // 配置key
	Value string `:"value"` // 配置value
}

type __diesel_schema_migrations struct {
	Version string `:"version"`
	RunOn   string `:"run_on"`
}
