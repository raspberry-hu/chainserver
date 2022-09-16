package response

type CollectionView struct {
	CollectionId   int     `json:"collection_id"`    // rankingID
	LogoImageURL   string  `json:"logo_image_url"`   // 集合标志图片
	BannerImageURL string  `json:"banner_image_url"` // 背景图片
	CollectionName string  `json:"collection_name"`  // 集合名称
	CollectionDesc string  `json:"collection_desc"`  // 集合描述
	CreateTax      int     `json:"create_tax"`       // 铸造税 1
	CategoryName   string  `json:"category_name"`    // 类别名称
	ChainName      string  `json:"chain_name"`       // 区块链名称
	CurrencyName   string  `json:"currency_name"`    // 区块链的币种ETH\BNB\AVAX
	Owner          int     `json:"owner"`            // 创建者
	OwnerName      string  `json:"owner_name"`       // 创建者名字
	Items          int     `json:"items"`            // 资产总量
	Favorites      int     `json:"favorites"`        // 收藏资产数额
	Amount         float32 `json:"amount"`           // 当前集合下的交易总额
}

//主页notable和trending
type HomeCollectionView struct {
	CollectionId   int    `json:"collection_id"`    // rankingID
	LogoImageURL   string `json:"logo_image_url"`   // 集合标志图片
	BannerImageURL string `json:"banner_image_url"` // 背景图片
	CollectionName string `json:"collection_name"`  // 集合名称
	CollectionDesc string `json:"collection_desc"`  // 集合描述
	CategoryName   string `json:"category_name"`    // 类别名称
	ChainName      string `json:"chain_name"`       // 区块链名称
	OwnerName      string `json:"wallet_addr"`      // 创建者地址
	Owner          int    `json:"owner"`            // 创建者id
}

//主页top collections和ranking
type HomeCollectionTop struct {
	CollectionId   int     `json:"collection_id"`   // rankingID
	LogoImageURL   string  `json:"logo_image_url"`  // 集合标志图片
	CollectionName string  `json:"collection_name"` // 集合名称
	CollectionDesc string  `json:"collection_desc"` // 集合描述
	CreateTax      int     `json:"create_tax"`      // 铸造税 1
	Amount         float32 `json:"amount"`          // 当前集合下的交易总额
	OwnerName      string  `json:"owner_name"`      // 拥有着的名字
	Owner          int     `json:"owner"`           // 拥有者id
}
