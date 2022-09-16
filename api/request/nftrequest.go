package request

// sn,nft_name,nft_desc,media_uri,create_tax,creater,owner
//type NftNew struct {
//	SN           string `json:"sn"`             // 流水号
//	NftName      string `json:"nft_name"`       // 名称
//	NftDesc      string `json:"nft_desc"`       // 描述
//	RightsRules  string `json:"rights_rules"`   // 权益规则
//	Creater      int    `json:"creater"`        // nft创建者
//	MediaUri     string `json:"media_uri"`      // 数字资产地址
//	CreateTax    int    `json:"create_tax"`     // 铸造税
//	Owner        int    `json:"owner"`          // nft拥有者
//	ChainName    string `json:"chain_name"`     // 链类型
//	MediaIpfsUri string `json:"media_ipfs_uri"` // 原数据ipfsurl
//	MetaDataUri  string `json:"meta_data_uri"`  // 原数据文件 ipfsurl
//	CurrencyName string `json:"currency_name"`  // 币种
//	Lazy         int    `json:"lazy"`           // 是否免收gas费 0 否 1 是
//	Donation     int    `json:"donation"`       // 是否捐赠 0 否 1 是
//	CollectionId int    `json:"collection_id"`  // 集合id
//	CategoriesId int    `json:"categories_id"`  // 类别id
//	ExploreUri   string `json:"explore_uri"`    // 音视频资产用于展示的图片地址
//	TxHash       string `json:"tx_hash"`        // 类别id
//}

type NftBuy struct {
	NftOwner string `json:"nft_owner"` //购买者
	Price    string `json:"price"`     //购买价格
	TokenId  string `json:"token_id"`  //购买id
	NftCount string `json:"nft_count"` //购买id中的编号·
}

type NftNew struct {
	SN           string `json:"sn"`             // 流水号
	NftName      string `json:"nft_name"`       // 名称 ===> tokenId
	NftDesc      string `json:"nft_desc"`       // 描述 ===> 拥有者
	RightsRules  string `json:"rights_rules"`   // 权益规则 ===> 数量
	Creater      int    `json:"creater"`        // nft创建者
	MediaUri     string `json:"media_uri"`      // 数字资产地址 ===> hash值
	CreateTax    int    `json:"create_tax"`     // 铸造税 ===> 版权税
	Owner        int    `json:"owner"`          // nft拥有者
	ChainName    string `json:"chain_name"`     // 链类型
	MediaIpfsUri string `json:"media_ipfs_uri"` // 原数据ipfsurl
	MetaDataUri  string `json:"meta_data_uri"`  // 原数据文件 ipfsurl ===> token url
	CurrencyName string `json:"currency_name"`  // 币种
	Lazy         int    `json:"lazy"`           // 是否免收gas费 0 否 1 是
	Donation     int    `json:"donation"`       // 是否捐赠 0 否 1 是
	CollectionId int    `json:"collection_id"`  // 集合id
	CategoriesId int    `json:"categories_id"`  // 类别id
	ExploreUri   string `json:"explore_uri"`    // 音视频资产用于展示的图片地址
	TxHash       string `json:"tx_hash"`        // 类别id
	AntTokenId   int    `json:"ant_token_id"`   // 蚂蚁链的tokenId
	AntCount     int    `json:"ant_count"`      // 蚂蚁链的铸造个数
	AntNftOwner  string `json:"ant_nft_owner"`  // 蚂蚁链的nft铸造者
	AntTokenUrl  string `json:"ant_token_url"`  // 蚂蚁链的tokenUrl
	AntTxHash    string `json:"ant_tx_hash"`    // 蚂蚁链藏品对应公链的hash值
	AntPrice     int    `json:"ant_price"`      // 蚂蚁链对应藏品价格
}

type NftOwner struct {
	NftId int `json:"nft_id"` // tokenid
	Owner int `json:"owner"`
	Buyer int `json:"buyer"`
}
